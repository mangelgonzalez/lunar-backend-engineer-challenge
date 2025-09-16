package domain

import "lunar-backend-engineer-challenge/pkg/domain"

const rocketNotExists = "Rocket not exists"

type RocketNotExists struct {
	domain.DomainError
	extraItems map[string]interface{}
}

func RocketNotExistsFromId(id RocketId) *RocketNotExists {
	return &RocketNotExists{extraItems: map[string]interface{}{
		"rocket_id": string(id),
	}}
}

func (i RocketNotExists) ExtraItems() map[string]interface{} {
	return i.extraItems
}

func (i RocketNotExists) Error() string {
	return rocketNotExists
}
