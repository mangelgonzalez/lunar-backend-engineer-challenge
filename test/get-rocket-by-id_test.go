package test

import (
	"fmt"
	"lunar-backend-engineer-challenge/internal/rockets/module/rocket/domain"
	"net/http"
	"testing"
)

func TestGetRocketById(t *testing.T) {
	setUp()

	rocketId := "aRocketId"
	rocket := CreateRandomRocketWithId(domain.RocketId(rocketId))
	givenTheFollowingRocketExist(t, rocket)

	response := executeUnauthenticatedJsonRequest(
		t,
		http.MethodGet,
		fmt.Sprintf("/rockets/%s", rocketId),
		nil,
		EmptyHeaders(),
	)

	expectedResponse := fmt.Sprintf(`{"data":{"type":"rocket","id":"%s","attributes":{"launch_speed":%d,"mission":"%s","type":"%s"}}}`, string(rocket.Id()), rocket.LaunchSpeed(), rocket.Mission(), rocket.Class())

	checkResponse(t, http.StatusOK, expectedResponse, response)
}

func TestGetRocketByIdAndNotExist(t *testing.T) {
	setUp()

	rocketId := "aRocketId"
	response := executeUnauthenticatedJsonRequest(
		t,
		http.MethodGet,
		fmt.Sprintf("/rockets/%s", rocketId),
		nil,
		EmptyHeaders(),
	)

	expectedResponse := `{"errors":[{"id":"<<PRESENCE>>","title":"Not Found","detail":"Rocket not exists","status":"404","code":"not_found"}]}`

	checkResponse(t, http.StatusNotFound, expectedResponse, response)
}
