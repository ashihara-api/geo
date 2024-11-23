package mysql

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/ashihara-api/core/interface/datasource/mysql"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"

	"github.com/ashihara-api/geo/core/domain/entity"
	"github.com/ashihara-api/geo/core/domain/repository"
)

const (
	tablePrefectures = "prefectures"
)

type (
	prefectureImpl struct {
		dbType DBType
		db     mysql.DB
		logger *slog.Logger
	}

	prefecture struct {
		Name      string     `db:"name"`
		Ruby      string     `db:"ruby"`
		Code      string     `db:"code"`
		CreatedAt *time.Time `db:"created_at"`
		UpdatedAt *time.Time `db:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at"`
	}
)

func ToPrefectureEntity(input *prefecture) (output *entity.Prefecture) {
	if input == nil {
		return nil
	}
	return &entity.Prefecture{
		Name: input.Name,
		Ruby: input.Ruby,
		Code: input.Code,
	}
}

func PrefectureReader(db *sql.DB, logger *slog.Logger) repository.Prefectures {
	return prefectures(db, Reader, logger)
}

func PrefectureWriter(db *sql.DB, logger *slog.Logger) repository.Prefectures {
	return prefectures(db, Writer, logger)
}

func prefectures(
	db *sql.DB,
	dbType DBType,
	logger *slog.Logger,
) repository.Prefectures {
	return &prefectureImpl{
		db:     sqlx.NewDb(db, mysql.DriverName),
		dbType: dbType,
		logger: logger,
	}
}

func (d *prefectureImpl) FindAll(
	ctx context.Context,
) (
	output []*entity.Prefecture,
	err error,
) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("*")
	sb.From(tablePrefectures)
	sb.Where(
		sb.IsNull("deleted_at"),
	)
	sb.OrderBy("code").Asc()
	q, args := sb.Build()
	d.logger.DebugContext(
		ctx,
		"prefectures.FindAll",
		slog.String("dbType", string(d.dbType)),
		slog.Group(
			"request",
			slog.String("query", q),
			slog.Any("args", args),
		),
	)
	var rs []*prefecture
	if err = d.db.SelectContext(ctx, &rs, q, args...); err != nil {
		d.logger.ErrorContext(
			ctx,
			"prefectures.FindAll",
			slog.String("dbType", string(d.dbType)),
			slog.String("action", "db.SelectContext"),
			slog.Group(
				"request",
				slog.String("query", q),
				slog.Any("args", args),
			),
			slog.Any("error", err),
		)
		if errors.Is(err, sql.ErrNoRows) {
			return []*entity.Prefecture{}, nil
		}
		return nil, err
	}
	if len(rs) == 0 {
		return []*entity.Prefecture{}, nil
	}
	output = make([]*entity.Prefecture, len(rs))
	for i, r := range rs {
		output[i] = ToPrefectureEntity(r)
	}
	return
}
