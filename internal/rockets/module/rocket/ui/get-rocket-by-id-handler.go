package ui

import (
	"errors"
	"lunar-backend-engineer-challenge/internal/rockets/module/rocket/application"
	"lunar-backend-engineer-challenge/internal/rockets/module/rocket/domain"
	"lunar-backend-engineer-challenge/pkg/bus/query"
	"lunar-backend-engineer-challenge/pkg/http/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func HandleGetRocketById(q query.Bus, jarm *middleware.JsonApiResponseMiddleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rocketId := mux.Vars(r)["rocketId"]

		if rocketId == "" {
			err := errors.New(`one or more of the required parameters are missing: ["rocketId"]`)
			jarm.WriteErrorResponse(w, middleware.BadRequestJsonApiHttpResponse(err.Error()), http.StatusBadRequest, err)
			return
		}

		response, err := q.Ask(r.Context(), &application.FindRocketByIdQuery{RocketId: rocketId})
		switch err.(type) {
		case nil:
			rocketResponse, ok := response.(*application.RocketResponse)
			if !ok {
				jarm.WriteErrorResponse(w, middleware.InternalServerErrorJsonApiHttpResponse(), http.StatusInternalServerError, err)
				return
			}
			jarm.WriteResponse(w, fromRocketResponse(rocketResponse), http.StatusOK)
			return
		case *domain.RocketNotExists:
			jarm.WriteErrorResponse(w, middleware.NotFoundRequestJsonApiHttpResponse(err.Error()), http.StatusNotFound, err)
			return
		default:
			jarm.WriteErrorResponse(w, middleware.InternalServerErrorJsonApiHttpResponse(), http.StatusInternalServerError, err)
			return
		}
	}
}

type RocketResponse struct {
	RocketId    string `jsonapi:"primary,rocket"`
	Type        string `jsonapi:"attr,type"`
	LaunchSpeed uint   `jsonapi:"attr,launch_speed"`
	Mission     string `jsonapi:"attr,mission"`
}

func fromRocketResponse(response *application.RocketResponse) *RocketResponse {
	return &RocketResponse{
		RocketId:    response.Id,
		Type:        response.Class,
		LaunchSpeed: response.LaunchSpeed,
		Mission:     response.Mission,
	}
}
