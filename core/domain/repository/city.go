package repository

import (
	"context"

	"github.com/ashihara-api/geo/core/domain/entity"
)

type (
	CitySearcher interface {
		SearchByPrefectureCode(ctx context.Context, code string) (output []*entity.City, err error)
	}

	CityCreater interface {
		Create(ctx context.Context, input []*entity.City) (err error)
	}

	Cities interface {
		CityCreater
		CitySearcher
	}

	CrawledData struct {
		CityCode       string
		CityName       string
		CityRuby       string
		PrefectureName string
		PrefectureRuby string
	}

	Crawler interface {
		Crawl(ctx context.Context) (output []*CrawledData, err error)
	}
)
