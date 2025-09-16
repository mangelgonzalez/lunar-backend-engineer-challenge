package ui

import (
	"encoding/json"
	"lunar-backend-engineer-challenge/internal/rockets/module/rocket/application"
	"lunar-backend-engineer-challenge/internal/rockets/module/rocket/domain"
	"lunar-backend-engineer-challenge/pkg/bus/command"
	"lunar-backend-engineer-challenge/pkg/http/middleware"
	"net/http"
)

func HandleStoreRocketMessages(c command.Bus, jarm *middleware.JsonApiResponseMiddleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := &StoreRocketRequest{}

		if err := json.NewDecoder(r.Body).Decode(request); err != nil {
			jarm.WriteErrorResponse(w, middleware.BadRequestJsonApiHttpResponse("Invalid request received"), http.StatusBadRequest, err)
			return
		}

		err := c.Dispatch(r.Context(), &application.StoreRocketMessagesCommand{
			RocketId:    request.MetaData.Channel,
			MessageType: request.MetaData.MessageType,
			Type:        request.Message.Type,
			LaunchSpeed: request.Message.LaunchSpeed,
			Mission:     request.Message.Mission,
			By:          request.Message.By,
			Reason:      request.Message.Reason,
			NewMission:  request.Message.NewMission,
		})

		switch err.(type) {
		case nil:
			jarm.WriteResponse(w, nil, http.StatusCreated)
			return
		case *domain.RocketNotExists:
			jarm.WriteErrorResponse(w, middleware.NotFoundRequestJsonApiHttpResponse(err.Error()), http.StatusNotFound, err)
		case *domain.RocketAlreadyExists:
			jarm.WriteErrorResponse(w, middleware.BadRequestJsonApiHttpResponse(err.Error()), http.StatusBadRequest, err)
			return
		default:
			jarm.WriteErrorResponse(w, middleware.InternalServerErrorJsonApiHttpResponse(), http.StatusInternalServerError, err)
			return
		}
	}
}

type StoreRocketRequest struct {
	MetaData struct {
		Channel       string `json:"channel"`
		MessageNumber int    `json:"messageNumber"`
		MessageTime   string `json:"messageTime"`
		MessageType   string `json:"messageType"`
	} `json:"metadata"`
	Message struct {
		Type        string `json:"type,omitempty"`
		LaunchSpeed uint   `json:"launchSpeed,omitempty"`
		Mission     string `json:"mission,omitempty"`
		By          uint   `json:"by,omitempty"`
		Reason      string `json:"reason,omitempty"`
		NewMission  string `json:"newMission,omitempty"`
	} `json:"message"`
}
