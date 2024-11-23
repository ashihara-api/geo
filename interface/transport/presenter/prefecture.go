package presenter

import (
	"context"

	"github.com/ashihara-api/geo/core/domain/usecase"
)

type (
	PrefectureBloc interface {
		FindAll(
			ctx context.Context,
		) (output *usecase.PrefectureAllFinderOutput, err error)
	}

	prefectureBlocImpl struct {
		findUsecase usecase.PrefectureAllFinder
	}
)

func NewPrefectureBloc(
	findUsecase usecase.PrefectureAllFinder,
) PrefectureBloc {
	return &prefectureBlocImpl{
		findUsecase: findUsecase,
	}
}

func (b *prefectureBlocImpl) FindAll(
	ctx context.Context,
) (output *usecase.PrefectureAllFinderOutput, err error) {
	return b.findUsecase.Execute(ctx)
}
