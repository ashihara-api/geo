package di

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/ashihara-api/geo/core/interactor/usecase"
	"github.com/ashihara-api/geo/interface/datasource/crawler"
	"github.com/ashihara-api/geo/interface/datasource/mysql"
	"github.com/ashihara-api/geo/interface/transport/presenter"
)

func City(c *http.Client, wdb, rdb *sql.DB, logger *slog.Logger) presenter.CityBloc {
	writer := mysql.CityWriter(wdb, logger)
	reader := mysql.CityReader(rdb, logger)
	client := crawler.Crawler(c, logger)

	prefectures := mysql.PrefectureReader(rdb, logger)

	return presenter.NewCityBloc(
		usecase.Import(client, writer, prefectures, logger),
		usecase.SearchCities(reader, logger),
	)
}
