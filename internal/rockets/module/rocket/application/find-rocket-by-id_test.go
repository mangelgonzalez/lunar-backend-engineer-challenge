package application_test

import (
	"context"
	"lunar-backend-engineer-challenge/internal/rockets/module/rocket/application"
	"lunar-backend-engineer-challenge/internal/rockets/module/rocket/domain"
	"lunar-backend-engineer-challenge/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindRocketByIdSuccess(t *testing.T) {
	rootCtx := context.Background()

	rocketId := domain.RocketId("aRocketId")
	t.Run("FindById Rocket by Id", func(t *testing.T) {
		repository := new(RocketRepositoryMock)
		handler := application.NewFindRocketByIdQueryHandler(repository)
		query := application.FindRocketByIdQuery{RocketId: string(rocketId)}

		expectedRocket := test.CreateRandomRocketWithId(rocketId)

		repository.ShouldFindById(rocketId, expectedRocket)

		rocket, err := handler.Handle(rootCtx, &query)

		assert.Equal(t, application.NewRocketResponse(*expectedRocket), rocket)
		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, repository)
	})
}

func TestFindRocketByIdWhenFail(t *testing.T) {
	rootCtx := context.Background()

	rocketId := domain.RocketId("aRocketId")
	t.Run("Rocket not found for given Id", func(t *testing.T) {
		repository := new(RocketRepositoryMock)
		handler := application.NewFindRocketByIdQueryHandler(repository)
		query := application.FindRocketByIdQuery{RocketId: string(rocketId)}

		repository.ShouldFindByIdAndReturnNotExist(rocketId)

		_, err := handler.Handle(rootCtx, &query)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "Rocket not exists")
		mock.AssertExpectationsForObjects(t, repository)
	})

	t.Run("Invalid DTO", func(t *testing.T) {
		repository := new(RocketRepositoryMock)
		handler := application.NewFindRocketByIdQueryHandler(repository)

		_, err := handler.Handle(rootCtx, &test.FakeDto{})

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "Invalid query")
		mock.AssertExpectationsForObjects(t, repository)
	})
}
