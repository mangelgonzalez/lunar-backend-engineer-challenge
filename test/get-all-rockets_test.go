package test

import (
	"fmt"
	"lunar-backend-engineer-challenge/internal/rockets/module/rocket/domain"
	"net/http"
	"testing"
)

func TestGetAllRockets(t *testing.T) {
	setUp()

	rocket := domain.NewRocketFromPrimitives("aRocketId", "Pam", 287, "ride")
	rocket2 := domain.NewRocketFromPrimitives("aRocket2Id", "Hobbit", 242, "play")
	givenTheFollowingRocketExist(t, rocket)
	givenTheFollowingRocketExist(t, rocket2)

	response := executeUnauthenticatedJsonRequest(
		t,
		http.MethodGet,
		"/rockets",
		nil,
		EmptyHeaders(),
	)

	expectedResponse := fmt.Sprintf(`{"data":[{"type":"rocket","id":"aRocket2Id","attributes":{"launch_speed":242,"mission":"play","type":"Hobbit"}},{"type":"rocket","id":"aRocketId","attributes":{"launch_speed":287,"mission":"ride","type":"Pam"}}]}`)

	checkResponse(t, http.StatusOK, expectedResponse, response)
}

func TestGetAllRocketsAndNotExist(t *testing.T) {
	setUp()

	response := executeUnauthenticatedJsonRequest(
		t,
		http.MethodGet,
		"/rockets",
		nil,
		EmptyHeaders(),
	)

	expectedResponse := `{"data":[]}`

	checkResponse(t, http.StatusOK, expectedResponse, response)
}
