package domain

import (
	"lunar-backend-engineer-challenge/pkg/utils"
)

type Rockets []*Rocket

type RocketId utils.Uuid

type Rocket struct {
	id          RocketId
	class       string
	launchSpeed uint
	mission     string
}

func NewRocket(id RocketId, class string, launchSpeed uint, mission string) *Rocket {
	return &Rocket{
		id:          id,
		class:       class,
		launchSpeed: launchSpeed,
		mission:     mission,
	}
}

func NewRocketFromPrimitives(id string, class string, launchSpeed uint, mission string) *Rocket {
	return &Rocket{
		id:          RocketId(id),
		class:       class,
		launchSpeed: launchSpeed,
		mission:     mission,
	}
}

func (r *Rocket) Id() RocketId {
	return r.id
}

func (r *Rocket) Class() string {
	return r.class
}

func (r *Rocket) LaunchSpeed() uint {
	return r.launchSpeed
}

func (r *Rocket) Mission() string {
	return r.mission
}

func (r *Rocket) IncreaseSpeed(speed uint) {
	r.launchSpeed += speed
}

func (r *Rocket) DecreaseSpeed(speed uint) {
	if r.launchSpeed-speed < 0 {
		r.launchSpeed = 0
	} else {
		r.launchSpeed -= speed
	}
}

func (r *Rocket) ChangeMission(newMission string) {
	r.mission = newMission
}
