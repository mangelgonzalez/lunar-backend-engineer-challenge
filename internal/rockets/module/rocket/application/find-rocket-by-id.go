package application

import (
	"context"
	"lunar-backend-engineer-challenge/internal/rockets/module/rocket/domain"
	"lunar-backend-engineer-challenge/pkg/bus"
)

type FindRocketByIdQuery struct {
	RocketId string
}

func (r FindRocketByIdQuery) Id() string { return "find_rocket_by_id" }

type FindRocketByIdQueryHandler struct {
	repository domain.RocketRepository
}

func NewFindRocketByIdQueryHandler(repository domain.RocketRepository) *FindRocketByIdQueryHandler {
	return &FindRocketByIdQueryHandler{repository: repository}
}

func (s *FindRocketByIdQueryHandler) Handle(ctx context.Context, query bus.Dto) (interface{}, error) {
	findQuery, ok := query.(*FindRocketByIdQuery)
	if !ok {
		return nil, bus.NewInvalidDto("Invalid query")
	}

	rocket, err := s.repository.FindById(domain.RocketId(findQuery.RocketId))
	if err != nil {
		return nil, err
	}
	return NewRocketResponse(*rocket), nil
}
