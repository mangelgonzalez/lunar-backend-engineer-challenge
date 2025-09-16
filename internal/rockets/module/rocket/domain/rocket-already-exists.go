package domain

import "lunar-backend-engineer-challenge/pkg/domain"

const rocketAlreadyExistsMessage = "Rocket already exists"

type RocketAlreadyExists struct {
	domain.DomainError
	extraItems map[string]interface{}
}

func NewRocketAlreadyExists(id RocketId) *RocketAlreadyExists {
	return &RocketAlreadyExists{
		extraItems: map[string]interface{}{
			"id": string(id),
		},
	}
}

func (i RocketAlreadyExists) ExtraItems() map[string]interface{} {
	return i.extraItems
}

func (i RocketAlreadyExists) Error() string {
	return rocketAlreadyExistsMessage
}
