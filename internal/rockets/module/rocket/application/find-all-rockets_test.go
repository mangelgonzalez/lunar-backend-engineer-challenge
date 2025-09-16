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

func TestFindAllRocketsSuccess(t *testing.T) {
	rootCtx := context.Background()

	rocketId := domain.RocketId("aRocketId")
	rocket2Id := domain.RocketId("aRocket2Id")
	t.Run("Find all Rockets", func(t *testing.T) {
		repository := new(RocketRepositoryMock)
		handler := application.NewFindAllRocketsQueryHandlerHandler(repository)
		query := application.FindAllRocketsQuery{}

		expectedRockets := domain.Rockets{test.CreateRandomRocketWithId(rocketId), test.CreateRandomRocketWithId(rocket2Id)}

		repository.ShouldFindAll(expectedRockets)

		rockets, err := handler.Handle(rootCtx, &query)

		assert.Equal(t, application.NewRocketsResponse(expectedRockets), rockets)
		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, repository)
	})
}

func TestFindAllRocketsWhenFail(t *testing.T) {
	rootCtx := context.Background()

	t.Run("Invalid DTO", func(t *testing.T) {
		repository := new(RocketRepositoryMock)
		handler := application.NewFindAllRocketsQueryHandlerHandler(repository)

		_, err := handler.Handle(rootCtx, &test.FakeDto{})

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "Invalid query")
		mock.AssertExpectationsForObjects(t, repository)
	})
}
