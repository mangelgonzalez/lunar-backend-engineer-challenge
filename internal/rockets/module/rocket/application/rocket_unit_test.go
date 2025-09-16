package application_test

import (
	"lunar-backend-engineer-challenge/internal/rockets/module/rocket/domain"

	"github.com/stretchr/testify/mock"
)

type RocketRepositoryMock struct {
	mock.Mock
}

func (r *RocketRepositoryMock) FindById(id domain.RocketId) (*domain.Rocket, error) {
	args := r.Called(id)

	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*domain.Rocket), nil
}

func (r *RocketRepositoryMock) FindAll() (domain.Rockets, error) {
	args := r.Called()

	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(domain.Rockets), nil
}
func (r *RocketRepositoryMock) Save(rocket *domain.Rocket) error {
	args := r.Called(rocket)
	return args.Error(0)
}

func (r *RocketRepositoryMock) DeleteById(id domain.RocketId) error {
	args := r.Called(id)
	return args.Error(0)
}

func (r *RocketRepositoryMock) ShouldFindById(id domain.RocketId, rocket *domain.Rocket) {
	r.On("FindById", id).Once().Return(rocket)
}

func (r *RocketRepositoryMock) ShouldFindByIdAndReturnNotExist(id domain.RocketId) {
	r.On("FindById", id).Once().Return(nil, domain.RocketNotExistsFromId(id))
}

func (r *RocketRepositoryMock) ShouldSaveRocket(rocket *domain.Rocket) {
	r.On("Save", rocket).Once().Return(nil)
}

func (r *RocketRepositoryMock) ShouldFindAll(rockets domain.Rockets) {
	r.On("FindAll").Once().Return(rockets)
}

func (r *RocketRepositoryMock) ShouldSaveRocketAndFails(rocket *domain.Rocket, err error) {
	r.On("Save", rocket).Once().Return(err)
}

func (r *RocketRepositoryMock) ShouldDeleteRocket(id domain.RocketId) {
	r.On("DeleteById", id).Once().Return(nil)
}

func (r *RocketRepositoryMock) ShouldDeleteRocketAndFails(id domain.RocketId, err error) {
	r.On("DeleteById", id).Once().Return(err)
}
