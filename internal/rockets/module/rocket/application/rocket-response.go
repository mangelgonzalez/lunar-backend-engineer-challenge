package application

import "lunar-backend-engineer-challenge/internal/rockets/module/rocket/domain"

type RocketResponse struct {
	Id          string
	Class       string
	LaunchSpeed uint
	Mission     string
}

func NewRocketResponse(rocket domain.Rocket) *RocketResponse {
	return &RocketResponse{
		Id:          string(rocket.Id()),
		Class:       rocket.Class(),
		LaunchSpeed: rocket.LaunchSpeed(),
		Mission:     rocket.Mission(),
	}
}

type RocketsResponse []*RocketResponse

func NewRocketsResponse(rockets domain.Rockets) RocketsResponse {
	response := make(RocketsResponse, len(rockets))

	for i, rocket := range rockets {
		response[i] = NewRocketResponse(*rocket)
	}

	return response
}
