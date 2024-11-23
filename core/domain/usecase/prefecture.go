package usecase

import (
	"context"

	"github.com/ashihara-api/geo/core/domain/entity"
)

type (
	PrefectureAllFinderOutput struct {
		Prefectures []*entity.Prefecture
	}

	PrefectureAllFinder interface {
		Execute(ctx context.Context) (output *PrefectureAllFinderOutput, err error)
	}
)
