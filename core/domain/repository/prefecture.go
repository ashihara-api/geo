package repository

import (
	"context"

	"github.com/ashihara-api/geo/core/domain/entity"
)

type (
	PrefectureAllFinder interface {
		FindAll(ctx context.Context) (output []*entity.Prefecture, err error)
	}

	Prefectures interface {
		PrefectureAllFinder
	}
)
