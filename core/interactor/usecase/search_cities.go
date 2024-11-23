package usecase

import (
	"context"
	"log/slog"

	"github.com/ashihara-api/core/domain/errors"

	"github.com/ashihara-api/geo/core/domain/repository"
	"github.com/ashihara-api/geo/core/domain/usecase"
)

type (
	implSearchCities struct {
		searcher repository.CitySearcher
		logger   *slog.Logger
	}
)

func SearchCities(
	searcher repository.CitySearcher,
	logger *slog.Logger,
) usecase.CitySearcher {
	return &implSearchCities{
		searcher: searcher,
		logger:   logger,
	}
}

func (u *implSearchCities) Execute(
	ctx context.Context,
	input usecase.CitySeacherInput,
) (
	output *usecase.CitySearcherOutput,
	err error,
) {
	rs, err := u.searcher.SearchByPrefectureCode(ctx, input.PrefectureCode)
	if err != nil {
		u.logger.ErrorContext(ctx,
			"SearchCities.Execute",
			slog.String("action", "searcher.FindByPrefectureCode"),
			slog.Any("error", err),
		)
		return nil, errors.NewCause(err, errors.CaseBackendError)
	}

	return &usecase.CitySearcherOutput{
		Cities: rs,
	}, nil
}
