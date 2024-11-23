package usecase

import (
	"context"

	"github.com/ashihara-api/core/domain"

	"github.com/ashihara-api/geo/core/domain/entity"
)

type (
	CitySeacherInput struct {
		PrefectureCode string
	}

	CitySearcherOutput struct {
		Cities []*entity.City
	}

	CitySearcher = domain.Usecase[CitySeacherInput, CitySearcherOutput]

	Importer interface {
		Execute(ctx context.Context) (err error)
	}
)
