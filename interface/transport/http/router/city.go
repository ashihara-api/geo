package router

import (
	"net/http"

	"github.com/ashihara-api/core/domain/errors"
	"github.com/ashihara-api/core/interface/transport/http/render"
	"github.com/ashihara-api/core/interface/transport/http/router"
	"github.com/go-chi/chi/v5"

	"github.com/ashihara-api/geo/core/domain/entity"
	"github.com/ashihara-api/geo/core/domain/usecase"
	"github.com/ashihara-api/geo/interface/transport/presenter"
)

type (
	// CityRouter ...
	CityRouter struct {
		bloc presenter.CityBloc
	}

	City struct {
		Name           string `json:"name"`
		Ruby           string `json:"ruby"`
		PrefectureCode string `json:"prefecture_code"`
		CityCode       string `json:"city_code"`
		CheckDigit     int    `json:"check_digit"`
	}

	CitySearchResponse struct {
		Cities []*City `json:"cities"`
	}
)

func fromCityEntity(input *entity.City) (output *City) {
	if input == nil {
		return nil
	}
	return &City{
		Name:           input.Name,
		Ruby:           input.Ruby,
		PrefectureCode: input.PrefectureCode,
		CityCode:       input.CityCode,
		CheckDigit:     input.CheckDigit,
	}
}

func fromCityEntities(input []*entity.City) (output []*City) {
	output = make([]*City, len(input))

	for i, v := range input {
		output[i] = fromCityEntity(v)
	}
	return
}

func toCitySearchResponse(input *usecase.CitySearcherOutput) (output *CitySearchResponse) {
	cities := fromCityEntities(input.Cities)
	return &CitySearchResponse{
		Cities: cities,
	}
}

func NewCity(bloc presenter.CityBloc) router.HTTPRouter {
	return &CityRouter{
		bloc: bloc,
	}
}

func (rt *CityRouter) Route(mux *chi.Mux) (err error) {
	routes := router.Route{
		Endpoints: []router.EndpointPattern{
			{
				Pattern: "/cities",
				Endpoints: map[string]router.Endpoint{
					http.MethodPost: {
						Handler: rt.Import,
					},
				},
			},
			{
				Pattern: "/cities",
				Endpoints: map[string]router.Endpoint{
					http.MethodGet: {
						Handler: rt.Search,
					},
				},
			},
		},
	}
	r := router.New(routes)
	return r.Build(mux)
}

func (rt *CityRouter) Import(w http.ResponseWriter, r *http.Request) {
	err := rt.bloc.Import(r.Context())
	if err != nil {
		render.ErrorJSON(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (rt *CityRouter) Search(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("prefecture_code")
	if code == "" {
		render.ErrorJSON(w, errors.NewCause(errors.New("prefecture code is empty"), errors.CaseBadRequest))
		return
	}

	output, err := rt.bloc.Search(r.Context(), usecase.CitySeacherInput{
		PrefectureCode: code,
	})
	if err != nil {
		render.ErrorJSON(w, err)
		return
	}
	render.JSON(w, http.StatusOK, toCitySearchResponse(output))
}
