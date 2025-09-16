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

func TestStoreRocketMessagesSuccess(t *testing.T) {
	rootCtx := context.Background()

	repository := new(RocketRepositoryMock)
	storer := domain.NewRocketStorer(repository)
	handler := application.NewStoreRocketMessagesCommandHandler(storer)

	rocketId := domain.RocketId("aRocketId")
	expectedRocket := test.CreateRandomRocketWithId(rocketId)

	t.Run("Store RocketLaunched message", func(t *testing.T) {

		command := application.StoreRocketMessagesCommand{
			RocketId:    string(rocketId),
			MessageType: "RocketLaunched",
			Type:        expectedRocket.Class(),
			LaunchSpeed: expectedRocket.LaunchSpeed(),
			Mission:     expectedRocket.Mission(),
			By:          0,
			Reason:      "",
			NewMission:  "",
		}

		repository.ShouldFindByIdAndReturnNotExist(rocketId)
		repository.ShouldSaveRocket(expectedRocket)

		err := handler.Handle(rootCtx, &command)

		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, repository)
	})

	t.Run("Store RocketSpeedIncreased message", func(t *testing.T) {

		command := application.StoreRocketMessagesCommand{
			RocketId:    string(rocketId),
			MessageType: "RocketSpeedIncreased",
			Type:        "",
			LaunchSpeed: 0,
			Mission:     "",
			By:          500,
			Reason:      "",
			NewMission:  "",
		}

		repository.ShouldFindById(rocketId, expectedRocket)
		repository.ShouldSaveRocket(expectedRocket)

		err := handler.Handle(rootCtx, &command)

		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, repository)
	})

	t.Run("Store RocketSpeedDecreased message", func(t *testing.T) {

		command := application.StoreRocketMessagesCommand{
			RocketId:    string(rocketId),
			MessageType: "RocketSpeedDecreased",
			Type:        "",
			LaunchSpeed: 0,
			Mission:     "",
			By:          500,
			Reason:      "",
			NewMission:  "",
		}

		repository.ShouldFindById(rocketId, expectedRocket)
		repository.ShouldSaveRocket(expectedRocket)

		err := handler.Handle(rootCtx, &command)

		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, repository)
	})

	t.Run("Store RocketExploded message", func(t *testing.T) {

		command := application.StoreRocketMessagesCommand{
			RocketId:    string(rocketId),
			MessageType: "RocketExploded",
			Type:        "",
			LaunchSpeed: 0,
			Mission:     "",
			By:          0,
			Reason:      "PRESSURE_VESSEL_FAILURE",
			NewMission:  "",
		}

		repository.ShouldFindById(rocketId, expectedRocket)
		repository.ShouldDeleteRocket(rocketId)

		err := handler.Handle(rootCtx, &command)

		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, repository)
	})

	t.Run("Store RocketMissionChanged message", func(t *testing.T) {

		command := application.StoreRocketMessagesCommand{
			RocketId:    string(rocketId),
			MessageType: "RocketMissionChanged",
			Type:        "",
			LaunchSpeed: 0,
			Mission:     "",
			By:          0,
			Reason:      "",
			NewMission:  "SHUTTLE_MIR",
		}

		repository.ShouldFindById(rocketId, expectedRocket)
		repository.ShouldSaveRocket(expectedRocket)

		err := handler.Handle(rootCtx, &command)

		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, repository)
	})
}

func TestStoreRocketMessagesWhenFail(t *testing.T) {
	rootCtx := context.Background()

	repository := new(RocketRepositoryMock)
	storer := domain.NewRocketStorer(repository)
	handler := application.NewStoreRocketMessagesCommandHandler(storer)

	rocketId := domain.RocketId("aRocketId")
	expectedRocket := test.CreateRandomRocketWithId(rocketId)

	t.Run("Store RocketLaunched message when already exists", func(t *testing.T) {

		command := application.StoreRocketMessagesCommand{
			RocketId:    string(rocketId),
			MessageType: "RocketLaunched",
			Type:        expectedRocket.Class(),
			LaunchSpeed: expectedRocket.LaunchSpeed(),
			Mission:     expectedRocket.Mission(),
			By:          0,
			Reason:      "",
			NewMission:  "",
		}

		repository.ShouldFindById(rocketId, expectedRocket)

		err := handler.Handle(rootCtx, &command)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "Rocket already exists")
		mock.AssertExpectationsForObjects(t, repository)
	})

	t.Run("Store RocketExploded message when not exists", func(t *testing.T) {

		command := application.StoreRocketMessagesCommand{
			RocketId:    string(rocketId),
			MessageType: "RocketExploded",
			Type:        "",
			LaunchSpeed: 0,
			Mission:     "",
			By:          0,
			Reason:      "PRESSURE_VESSEL_FAILURE",
			NewMission:  "",
		}

		repository.ShouldFindByIdAndReturnNotExist(rocketId)

		err := handler.Handle(rootCtx, &command)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "Rocket not exists")
		mock.AssertExpectationsForObjects(t, repository)
	})

	t.Run("Invalid DTO", func(t *testing.T) {

		err := handler.Handle(rootCtx, &test.FakeDto{})

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "Invalid command")
		mock.AssertExpectationsForObjects(t, repository)
	})
}
