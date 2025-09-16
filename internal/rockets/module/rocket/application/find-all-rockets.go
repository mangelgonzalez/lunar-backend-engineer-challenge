package application

import (
	"context"
	"lunar-backend-engineer-challenge/internal/rockets/module/rocket/domain"
	"lunar-backend-engineer-challenge/pkg/bus"
)

type FindAllRocketsQuery struct{}

func (s FindAllRocketsQuery) Id() string { return "find_all_rockets" }

type FindAllRocketsQueryHandler struct {
	repository domain.RocketRepository
}

func NewFindAllRocketsQueryHandlerHandler(repository domain.RocketRepository) *FindAllRocketsQueryHandler {
	return &FindAllRocketsQueryHandler{repository: repository}
}

func (s *FindAllRocketsQueryHandler) Handle(_ context.Context, query bus.Dto) (interface{}, error) {
	_, ok := query.(*FindAllRocketsQuery)
	if !ok {
		return nil, bus.NewInvalidDto("Invalid query")
	}

	rockets, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	return NewRocketsResponse(rockets), nil
}
