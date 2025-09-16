package test

import (
	"lunar-backend-engineer-challenge/internal/rockets/module/rocket/domain"

	"github.com/brianvoe/gofakeit/v6"
)

func CreateRandomRocketWithId(id domain.RocketId) *domain.Rocket {
	return domain.NewRocketFromPrimitives(
		string(id),
		gofakeit.PetName(),
		uint(gofakeit.Number(1, 500)),
		gofakeit.VerbAction(),
	)
}

type FakeDto struct{}

func (f *FakeDto) Id() string { return "fake" }
