package domain

import "errors"

type RocketStorer struct {
	repository RocketRepository
}

func NewRocketStorer(repository RocketRepository) *RocketStorer {
	return &RocketStorer{repository: repository}
}

func (s *RocketStorer) Store(id RocketId, messageType string, class string, launchSpeed uint, mission string, by uint, newMission string) error {

	rocket, err := s.repository.FindById(id)
	if err != nil && !errors.As(err, &RocketNotExists{}) && messageType != "RocketLaunched" {
		return err
	}

	switch messageType {
	case "RocketExploded":
		err := s.repository.DeleteById(id)
		if err != nil {
			return err
		}

		return nil
	case "RocketLaunched":
		if rocket != nil && rocket.Id() == id {
			return NewRocketAlreadyExists(rocket.Id())
		}

		rocket = NewRocket(id, class, launchSpeed, mission)
	case "RocketSpeedIncreased":
		rocket.IncreaseSpeed(by)
	case "RocketSpeedDecreased":
		rocket.DecreaseSpeed(by)
	case "RocketMissionChanged":
		rocket.ChangeMission(newMission)
	}

	return s.repository.Save(rocket)
}
