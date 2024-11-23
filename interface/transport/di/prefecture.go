package di

import (
	"database/sql"
	"log/slog"

	"github.com/ashihara-api/geo/core/interactor/usecase"
	"github.com/ashihara-api/geo/interface/datasource/mysql"
	"github.com/ashihara-api/geo/interface/transport/presenter"
)

func Prefecture(rdb *sql.DB, logger *slog.Logger) presenter.PrefectureBloc {
	reader := mysql.PrefectureReader(rdb, logger)

	return presenter.NewPrefectureBloc(
		usecase.FindAllPrefectures(reader, logger),
	)
}
