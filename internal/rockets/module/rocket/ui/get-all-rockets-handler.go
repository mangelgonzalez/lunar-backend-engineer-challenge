package ui

import (
	"lunar-backend-engineer-challenge/internal/rockets/module/rocket/application"
	"lunar-backend-engineer-challenge/pkg/bus/query"
	"lunar-backend-engineer-challenge/pkg/http/middleware"
	"net/http"
)

func HandleGetAllRockets(q query.Bus, jarm *middleware.JsonApiResponseMiddleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		response, err := q.Ask(r.Context(), &application.FindAllRocketsQuery{})
		switch err.(type) {
		case nil:
			rocketsResponse, ok := response.(application.RocketsResponse)
			if !ok {
				jarm.WriteErrorResponse(w, middleware.InternalServerErrorJsonApiHttpResponse(), http.StatusInternalServerError, err)
				return
			}
			jarm.WriteResponse(w, fromRocketsResponse(rocketsResponse), http.StatusOK)
			return
		default:
			jarm.WriteErrorResponse(w, middleware.InternalServerErrorJsonApiHttpResponse(), http.StatusInternalServerError, err)
			return
		}
	}
}

type RocketsResponse []*RocketResponse

func fromRocketsResponse(rockets application.RocketsResponse) RocketsResponse {
	response := make(RocketsResponse, len(rockets))

	for i, rocket := range rockets {
		response[i] = &RocketResponse{
			RocketId:    rocket.Id,
			Type:        rocket.Class,
			LaunchSpeed: rocket.LaunchSpeed,
			Mission:     rocket.Mission,
		}
	}

	return response
}
