package usecase

import (
	"context"
	"log/slog"

	"github.com/ashihara-api/core/domain/errors"

	"github.com/ashihara-api/geo/core/domain/repository"
	"github.com/ashihara-api/geo/core/domain/usecase"
)

type (
	implFindAllPrefectures struct {
		finder repository.PrefectureAllFinder
		logger *slog.Logger
	}
)

func FindAllPrefectures(
	finder repository.PrefectureAllFinder,
	logger *slog.Logger,
) usecase.PrefectureAllFinder {
	return &implFindAllPrefectures{
		finder: finder,
		logger: logger,
	}
}

func (u *implFindAllPrefectures) Execute(
	ctx context.Context,
) (
	output *usecase.PrefectureAllFinderOutput,
	err error,
) {
	rs, err := u.finder.FindAll(ctx)
	if err != nil {
		u.logger.ErrorContext(ctx,
			"FindAllPrefectures.Execute",
			slog.String("action", "finder.FindAll"),
			slog.Any("error", err),
		)
		return nil, errors.NewCause(err, errors.CaseBackendError)
	}
	return &usecase.PrefectureAllFinderOutput{
		Prefectures: rs,
	}, nil
}
