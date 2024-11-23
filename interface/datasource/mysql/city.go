package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ashihara-api/core/interface/datasource/mysql"
	"github.com/ashihara-api/core/utils/jst"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"

	"github.com/ashihara-api/geo/core/domain/entity"
	"github.com/ashihara-api/geo/core/domain/repository"
)

const (
	tableCities = "cities"
)

type (
	cityImpl struct {
		dbType DBType
		db     mysql.DB
		logger *slog.Logger
	}

	city struct {
		Name           string     `db:"name"`
		Ruby           string     `db:"ruby"`
		PrefectureCode string     `db:"prefecture_code"`
		CityCode       string     `db:"city_code"`
		CheckDigit     int        `db:"check_digit"`
		CreatedAt      *time.Time `db:"created_at"`
		UpdatedAt      *time.Time `db:"updated_at"`
		DeletedAt      *time.Time `db:"deleted_at"`
	}
)

func ToCityEntity(input *city) (output *entity.City) {
	if input == nil {
		return nil
	}
	return &entity.City{
		Name:           input.Name,
		Ruby:           input.Ruby,
		PrefectureCode: input.PrefectureCode,
		CityCode:       input.CityCode,
		CheckDigit:     input.CheckDigit,
	}
}

func CityReader(db *sql.DB, logger *slog.Logger) repository.Cities {
	return cities(db, Reader, logger)
}

func CityWriter(db *sql.DB, logger *slog.Logger) repository.Cities {
	return cities(db, Writer, logger)
}

func cities(
	db *sql.DB,
	dbType DBType,
	logger *slog.Logger,
) repository.Cities {
	return &cityImpl{
		db:     sqlx.NewDb(db, mysql.DriverName),
		dbType: dbType,
		logger: logger,
	}
}

func (d *cityImpl) Create(
	ctx context.Context,
	input []*entity.City,
) (
	err error,
) {
	if d.dbType != Writer {
		err = repository.ErrNoPermission
		d.logger.ErrorContext(
			ctx,
			"cities.Create",
			slog.String("dbType", string(d.dbType)),
			slog.Any("error", err),
		)
		return
	}
	if input == nil {
		return
	}
	tt := jst.Now()
	_, err = transaction(d.db, func(tx *sqlx.Tx) (sql.Result, error) {
		for _, v := range input {
			c := city{
				Name:           v.Name,
				Ruby:           v.Ruby,
				PrefectureCode: v.PrefectureCode,
				CityCode:       v.CityCode,
				CheckDigit:     v.CheckDigit,
				CreatedAt:      &tt,
				UpdatedAt:      &tt,
			}
			_, err := d.insert(ctx, tx, &c)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	return
}

func (d *cityImpl) SearchByPrefectureCode(
	ctx context.Context,
	code string,
) (
	output []*entity.City,
	err error,
) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("*")
	sb.From(tableCities)
	sb.Where(
		sb.E("prefecture_code", code),
		sb.IsNull("deleted_at"),
	)
	sb.OrderBy("prefecture_code", "city_code").Asc()
	q, args := sb.Build()
	d.logger.DebugContext(
		ctx,
		"cities.FindByPrefectureCode",
		slog.String("dbType", string(d.dbType)),
		slog.Group(
			"request",
			slog.String("query", q),
			slog.Any("args", args),
		),
	)
	var rs []*city
	if err = d.db.SelectContext(ctx, &rs, q, args...); err != nil {
		d.logger.ErrorContext(
			ctx,
			"cities.FindByPrefectureCode",
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
			return []*entity.City{}, nil
		}
		return nil, err
	}
	if len(rs) == 0 {
		return []*entity.City{}, nil
	}
	output = make([]*entity.City, len(rs))
	for i, r := range rs {
		output[i] = ToCityEntity(r)
	}
	return
}

func (d *cityImpl) insert(
	ctx context.Context,
	tx *sqlx.Tx,
	r ...interface{},
) (
	result sql.Result,
	err error,
) {
	ps := sqlbuilder.NewStruct(new(city))
	ib := ps.InsertInto(tableCities, r...)
	q, args := ib.Build()
	d.logger.DebugContext(
		ctx,
		"cities.insert",
		slog.String("dbType", string(d.dbType)),
		slog.Group(
			"request",
			slog.String("query", q),
			slog.Any("args", args),
		),
	)
	if result, err = tx.ExecContext(ctx, q, args...); err != nil {
		d.logger.ErrorContext(
			ctx,
			"cities.insert",
			slog.String("dbType", string(d.dbType)),
			slog.String("action", "db.ExecContext"),
			slog.Group(
				"request",
				slog.String("query", q),
				slog.Any("args", args),
			),
			slog.Any("error", err),
		)
	}
	return
}

// transaction ...
func transaction(q sqlx.Queryer, txFunc func(*sqlx.Tx) (sql.Result, error)) (result sql.Result, err error) {
	switch db := q.(type) {
	case *sqlx.DB:
		var tx *sqlx.Tx
		tx, err = db.Beginx()
		if err != nil {
			return nil, err
		}

		defer func() {
			if p := recover(); p != nil {
				// recover
				_ = tx.Rollback()
				switch x := p.(type) {
				case error:
					err = x
				case string:
					err = errors.New(x)
				default:
					err = fmt.Errorf("unknown error: %v", x)
				}
			} else if err != nil {
				// rollback
				_ = tx.Rollback()
			} else {
				// commit
				err = tx.Commit()
			}
		}()
		result, err = txFunc(tx)
		return
	case *sqlx.Tx:
		defer func() {
			if p := recover(); p != nil {
				// recover
				_ = db.Rollback()
				switch x := p.(type) {
				case error:
					err = x
				case string:
					err = errors.New(x)
				default:
					err = fmt.Errorf("unknown error: %v", x)
				}
			} else if err != nil {
				// rollback
				_ = db.Rollback()
			} else {
				// commit
				err = db.Commit()
			}
		}()
		result, err = txFunc(db)
		return
	}
	return nil, errors.New("fail to start transaction. invalid request type")
}
