package crawler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/xuri/excelize/v2"
	"golang.org/x/text/width"

	"github.com/ashihara-api/geo/core/domain/repository"
)

type (
	crawlerImpl struct {
		client *http.Client
		logger *slog.Logger
	}
)

const (
	codeXlsxURL = "https://nlftp.mlit.go.jp/ksj/gml/codelist/AdminiBoundary_CD.xlsx"
)

func Crawler(c *http.Client, logger *slog.Logger) repository.Crawler {
	return &crawlerImpl{
		client: c,
		logger: logger,
	}
}

func (d *crawlerImpl) Crawl(ctx context.Context) (
	output []*repository.CrawledData,
	err error,
) {
	res, err := d.request(ctx)
	if err != nil {
		return nil, err
	}
	rb := res.Body
	defer func() {
		_ = rb.Close()
	}()

	f, err := excelize.OpenReader(rb)
	if err != nil {
		d.logger.ErrorContext(
			ctx,
			"crawler.Crawl",
			slog.String("action", "excelize.OpenReader"),
			slog.Any("error", err),
		)
		return nil, err
	}
	rows, err := f.GetRows("行政区域コード")
	for rowNum, row := range rows {
		if rowNum <= 3 {
			// rowNum <= 3: not code table
			// rowNum = 3: headers
			continue
		}
		if len(row) < 3 || row[2] == "" {
			// (市区町村名 （漢字）is empty) == 都道府県
			continue
		}
		var v repository.CrawledData
		ignore := false
		for colNum, cell := range row {
			if colNum >= 6 {
				break
			}
			switch colNum {
			case 0:
				// 行政区域コード
				v.CityCode = cell
			case 1:
				// 都道府県名（漢字）
				v.PrefectureName = cell
			case 2:
				// 市区町村名（漢字）
				v.CityName = cell
			case 3:
				// 都道府県名（ｶﾅ）
				// hankaku -> zenkaku
				v.PrefectureRuby = width.Widen.String(cell)
			case 4:
				// 市区町村名（ｶﾅ）
				// hankaku -> zenkaku
				v.CityRuby = width.Widen.String(cell)
			case 5:
				if cell != "" {
					ignore = true
				}
			}
		}
		if !ignore {
			output = append(output, &v)
		}
	}
	return
}

func (d *crawlerImpl) request(
	ctx context.Context,
) (res *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, codeXlsxURL, http.NoBody)
	if err != nil {
		d.logger.ErrorContext(
			ctx,
			"crawler.request",
			slog.String("action", "http.NewRequestWithContext"),
			slog.Any("error", err),
		)
		return
	}

	res, err = d.client.Do(req)
	if err != nil {
		d.logger.ErrorContext(
			ctx,
			"crawler.request",
			slog.String("action", "client.Do"),
			slog.Any("error", err),
		)
		return
	}
	return
}
