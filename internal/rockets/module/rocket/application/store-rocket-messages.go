package application

import (
	"context"
	"lunar-backend-engineer-challenge/internal/rockets/module/rocket/domain"
	"lunar-backend-engineer-challenge/pkg/bus"
)

type StoreRocketMessagesCommand struct {
	RocketId    string
	MessageType string
	Type        string
	LaunchSpeed uint
	Mission     string
	By          uint
	Reason      string
	NewMission  string
}

func (s StoreRocketMessagesCommand) Id() string { return "store_rocket_messages_command" }

type StoreRocketMessagesCommandHandler struct {
	storer *domain.RocketStorer
}

func NewStoreRocketMessagesCommandHandler(storer *domain.RocketStorer) *StoreRocketMessagesCommandHandler {
	return &StoreRocketMessagesCommandHandler{storer: storer}
}

func (s *StoreRocketMessagesCommandHandler) Handle(_ context.Context, command bus.Dto) error {
	rocketsCommand, ok := command.(*StoreRocketMessagesCommand)
	if !ok {
		return bus.NewInvalidDto("Invalid command")
	}

	return s.storer.Store(
		domain.RocketId(rocketsCommand.RocketId),
		rocketsCommand.MessageType,
		rocketsCommand.Type,
		rocketsCommand.LaunchSpeed,
		rocketsCommand.Mission,
		rocketsCommand.By,
		rocketsCommand.NewMission,
	)
}
