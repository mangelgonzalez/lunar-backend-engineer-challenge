package test

import (
	"fmt"
	"lunar-backend-engineer-challenge/internal/rockets/module/rocket/domain"
	"net/http"
	"testing"
)

func TestStoreRocketLaunched(t *testing.T) {
	setUp()

	response := executeUnauthenticatedJsonRequest(
		t,
		http.MethodPost,
		"/rockets",
		[]byte(fmt.Sprintf(`{
    		"metadata": {
        		"channel": "193270a9-c9cf-404a-8f83-838e71d9ae67",
        		"messageNumber": 1,
				"messageTime": "2022-02-02T19:39:05.86337+01:00",
        		"messageType": "RocketLaunched"
    		},
    		"message": {
        		"type": "Falcon-9",
        		"launchSpeed": 500,
        		"mission": "ARTEMIS"
			}
		}`)),
		EmptyHeaders(),
	)

	checkResponse(t, http.StatusCreated, "", response)
}

func TestStoreRocketLaunchedAlreadyExist(t *testing.T) {
	setUp()

	rocket := domain.NewRocketFromPrimitives("aRocketId", "Falcon-9", 500, "ARTEMIS")
	givenTheFollowingRocketExist(t, rocket)

	response := executeUnauthenticatedJsonRequest(
		t,
		http.MethodPost,
		"/rockets",
		[]byte(fmt.Sprintf(`{
    		"metadata": {
        		"channel": "%s",
        		"messageNumber": 1,
				"messageTime": "2022-02-02T19:39:05.86337+01:00",
        		"messageType": "RocketLaunched"
    		},
    		"message": {
        		"type": "%s",
        		"launchSpeed": %d,
        		"mission": "%s"
			}
		}`, rocket.Id(), rocket.Class(), rocket.LaunchSpeed(), rocket.Mission())),
		EmptyHeaders(),
	)

	expectedResponse := `{"errors":[{"id":"<<PRESENCE>>","title":"Bad Request","detail":"Rocket already exists","status":"400","code":"bad_request"}]}`

	checkResponse(t, http.StatusBadRequest, expectedResponse, response)
}

func TestStoreRocketSpeedIncreased(t *testing.T) {
	setUp()

	rocket := domain.NewRocketFromPrimitives("aRocketId", "Falcon-9", 500, "ARTEMIS")
	givenTheFollowingRocketExist(t, rocket)

	response := executeUnauthenticatedJsonRequest(
		t,
		http.MethodPost,
		"/rockets",
		[]byte(fmt.Sprintf(`{
    		"metadata": {
        		"channel": "%s",
        		"messageNumber": 1,
				"messageTime": "2022-02-02T19:39:05.86337+01:00",
        		"messageType": "RocketSpeedIncreased"
    		},
    		"message": {
        		"by": 100
			}
		}`, rocket.Id())),
		EmptyHeaders(),
	)

	checkResponse(t, http.StatusCreated, "", response)
}

func TestStoreRocketSpeedIncreasedNotExist(t *testing.T) {
	setUp()

	rocketId := "aRocketId"

	response := executeUnauthenticatedJsonRequest(
		t,
		http.MethodPost,
		"/rockets",
		[]byte(fmt.Sprintf(`{
    		"metadata": {
        		"channel": "%s",
        		"messageNumber": 1,
				"messageTime": "2022-02-02T19:39:05.86337+01:00",
        		"messageType": "RocketSpeedIncreased"
    		},
    		"message": {
        		"by": 100
			}
		}`, rocketId)),
		EmptyHeaders(),
	)

	expectedResponse := `{"errors":[{"id":"<<PRESENCE>>","title":"Not Found","detail":"Rocket not exists","status":"404","code":"not_found"}]}`

	checkResponse(t, http.StatusNotFound, expectedResponse, response)
}

func TestStoreRocketSpeedDecreased(t *testing.T) {
	setUp()

	rocket := domain.NewRocketFromPrimitives("aRocketId", "Falcon-9", 500, "ARTEMIS")
	givenTheFollowingRocketExist(t, rocket)

	response := executeUnauthenticatedJsonRequest(
		t,
		http.MethodPost,
		"/rockets",
		[]byte(fmt.Sprintf(`{
    		"metadata": {
        		"channel": "%s",
        		"messageNumber": 1,
				"messageTime": "2022-02-02T19:39:05.86337+01:00",
        		"messageType": "RocketSpeedDecreased"
    		},
    		"message": {
        		"by": 100
			}
		}`, rocket.Id())),
		EmptyHeaders(),
	)

	checkResponse(t, http.StatusCreated, "", response)
}

func TestStoreRocketMissionChanged(t *testing.T) {
	setUp()

	rocket := domain.NewRocketFromPrimitives("aRocketId", "Falcon-9", 500, "ARTEMIS")
	givenTheFollowingRocketExist(t, rocket)

	response := executeUnauthenticatedJsonRequest(
		t,
		http.MethodPost,
		"/rockets",
		[]byte(fmt.Sprintf(`{
    		"metadata": {
        		"channel": "%s",
        		"messageNumber": 1,
				"messageTime": "2022-02-02T19:39:05.86337+01:00",
        		"messageType": "RocketMissionChanged"
    		},
    		"message": {
        		"newMission": "SHUTTLE_MIR"
			}
		}`, rocket.Id())),
		EmptyHeaders(),
	)

	checkResponse(t, http.StatusCreated, "", response)
}

func TestStoreRocketExploded(t *testing.T) {
	setUp()

	rocket := domain.NewRocketFromPrimitives("aRocketId", "Falcon-9", 500, "ARTEMIS")
	givenTheFollowingRocketExist(t, rocket)

	response := executeUnauthenticatedJsonRequest(
		t,
		http.MethodPost,
		"/rockets",
		[]byte(fmt.Sprintf(`{
    		"metadata": {
        		"channel": "%s",
        		"messageNumber": 1,
				"messageTime": "2022-02-02T19:39:05.86337+01:00",
        		"messageType": "RocketExploded"
    		},
    		"message": {
        		"reason": "PRESSURE_VESSEL_FAILURE"
			}
		}`, rocket.Id())),
		EmptyHeaders(),
	)

	checkResponse(t, http.StatusCreated, "", response)
}
