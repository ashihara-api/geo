package router

import (
	"net/http"

	"github.com/ashihara-api/core/interface/transport/http/render"
	"github.com/ashihara-api/core/interface/transport/http/router"
	"github.com/go-chi/chi/v5"

	"github.com/ashihara-api/geo/core/domain/entity"
	"github.com/ashihara-api/geo/core/domain/usecase"
	"github.com/ashihara-api/geo/interface/transport/presenter"
)

type (
	// PrefectureRouter ...
	PrefectureRouter struct {
		bloc presenter.PrefectureBloc
	}

	Prefecture struct {
		Code string `json:"code"`
		Name string `json:"name"`
		Ruby string `json:"ruby"`
	}

	PrefectureFindAllResponse struct {
		Prefectures []*Prefecture `json:"prefectures"`
	}
)

func fromPrefectureEntity(input *entity.Prefecture) (output *Prefecture) {
	if input == nil {
		return nil
	}
	return &Prefecture{
		Code: input.Code,
		Name: input.Name,
		Ruby: input.Ruby,
	}
}

func fromPrefectureEntities(input []*entity.Prefecture) (output []*Prefecture) {
	output = make([]*Prefecture, len(input))

	for i, v := range input {
		output[i] = fromPrefectureEntity(v)
	}
	return
}

func toPrefectureFindAllResponse(input *usecase.PrefectureAllFinderOutput) (output *PrefectureFindAllResponse) {
	prefectures := fromPrefectureEntities(input.Prefectures)
	return &PrefectureFindAllResponse{
		Prefectures: prefectures,
	}
}

func NewPrefecture(bloc presenter.PrefectureBloc) router.HTTPRouter {
	return &PrefectureRouter{
		bloc: bloc,
	}
}

func (rt *PrefectureRouter) Route(mux *chi.Mux) (err error) {
	routes := router.Route{
		Endpoints: []router.EndpointPattern{
			{
				Pattern: "/prefectures",
				Endpoints: map[string]router.Endpoint{
					http.MethodGet: {
						Handler: rt.FindAll,
					},
				},
			},
		},
	}
	r := router.New(routes)
	return r.Build(mux)
}

func (rt *PrefectureRouter) FindAll(w http.ResponseWriter, r *http.Request) {
	output, err := rt.bloc.FindAll(r.Context())
	if err != nil {
		render.ErrorJSON(w, err)
		return
	}
	render.JSON(w, http.StatusOK, toPrefectureFindAllResponse(output))
}
