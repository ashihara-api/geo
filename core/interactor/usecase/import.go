package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"

	"github.com/ashihara-api/core/domain/errors"

	"github.com/ashihara-api/geo/core/domain/entity"
	"github.com/ashihara-api/geo/core/domain/repository"
	"github.com/ashihara-api/geo/core/domain/usecase"
)

var (
	regCityCode = regexp.MustCompile(`^\d{5}$`)
)

type (
	implImport struct {
		crawler    repository.Crawler
		cities     repository.CityCreater
		prefecures repository.PrefectureAllFinder
		logger     *slog.Logger
	}
)

func generateCheckDigit(s string) (output int, err error) {
	if len(s) != 5 {
		return 0, fmt.Errorf("%s is invalid code", s)
	}
	if !regCityCode.MatchString(s) {
		return 0, fmt.Errorf("%s is invalid code", s)
	}
	a, _ := strconv.Atoi(string(s[0]))
	b, _ := strconv.Atoi(string(s[1]))
	c, _ := strconv.Atoi(string(s[2]))
	d, _ := strconv.Atoi(string(s[3]))
	e, _ := strconv.Atoi(string(s[4]))

	return (11 - (((a * 6) + (b * 5) + (c * 4) + (d * 3) + (e * 2)) % 11)) % 10, nil
}

func Import(
	crawler repository.Crawler,
	cities repository.CityCreater,
	prefecures repository.PrefectureAllFinder,
	logger *slog.Logger,
) usecase.Importer {
	return &implImport{
		crawler:    crawler,
		cities:     cities,
		prefecures: prefecures,
		logger:     logger,
	}
}

func (u *implImport) Execute(ctx context.Context) (err error) {
	rs, err := u.crawler.Crawl(ctx)
	if err != nil {
		u.logger.ErrorContext(ctx,
			"Import.Execute",
			slog.String("action", "crawler.Crawl"),
			slog.Any("error", err),
		)
		return errors.NewCause(err, errors.CaseBackendError)
	}
	ps, err := u.prefecures.FindAll(ctx)
	if err != nil {
		u.logger.ErrorContext(ctx,
			"Import.Execute",
			slog.String("action", "prefecures.FindAll"),
			slog.Any("error", err),
		)
		return errors.NewCause(err, errors.CaseBackendError)
	}
	// prefecureName: prefectureCode
	mp := map[string]string{}
	for _, p := range ps {
		mp[p.Name] = p.Code
	}
	cs := make([]*entity.City, len(rs))
	for i, r := range rs {
		pCode, ok := mp[r.PrefectureName]
		if !ok {
			err = fmt.Errorf("%s is a invalid prefecture name", r.PrefectureName)
			u.logger.ErrorContext(ctx,
				"Import.Execute",
				slog.String("action", "mp[r.PrefectureName]"),
				slog.Any("error", err),
			)
			return errors.NewCause(err, errors.CaseBackendError)
		}

		cd, err := generateCheckDigit(r.CityCode)
		if err != nil {
			u.logger.ErrorContext(ctx,
				"Import.Execute",
				slog.String("action", "generateCheckDigit"),
				slog.Any("error", err),
			)
			return errors.NewCause(err, errors.CaseBackendError)
		}

		cs[i] = &entity.City{
			Name:           r.CityName,
			Ruby:           r.CityRuby,
			PrefectureCode: pCode,
			CityCode:       r.CityCode,
			CheckDigit:     cd,
		}
	}
	if err = u.cities.Create(ctx, cs); err != nil {
		u.logger.ErrorContext(ctx,
			"Import.Execute",
			slog.String("action", "cities.Create"),
			slog.Any("error", err),
		)
		return errors.NewCause(err, errors.CaseBackendError)
	}
	return
}
