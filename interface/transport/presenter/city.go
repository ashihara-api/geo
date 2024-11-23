package presenter

import (
	"context"

	"github.com/ashihara-api/geo/core/domain/usecase"
)

type (
	CityBloc interface {
		Import(ctx context.Context) (err error)
		Search(
			ctx context.Context,
			input usecase.CitySeacherInput,
		) (output *usecase.CitySearcherOutput, err error)
	}

	cityBlocImpl struct {
		importUsecase usecase.Importer
		searchUsecase usecase.CitySearcher
	}
)

func NewCityBloc(
	importUsecase usecase.Importer,
	searchUsecase usecase.CitySearcher,
) CityBloc {
	return &cityBlocImpl{
		importUsecase: importUsecase,
		searchUsecase: searchUsecase,
	}
}

func (b *cityBlocImpl) Import(
	ctx context.Context,
) (err error) {
	return b.importUsecase.Execute(ctx)
}

func (b *cityBlocImpl) Search(
	ctx context.Context,
	input usecase.CitySeacherInput,
) (output *usecase.CitySearcherOutput, err error) {
	return b.searchUsecase.Execute(ctx, input)
}
