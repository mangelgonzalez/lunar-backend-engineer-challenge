package domain

type RocketRepository interface {
	RocketRepositoryReader
	RocketRepositoryWriter
}

type RocketRepositoryReader interface {
	FindById(id RocketId) (*Rocket, error)
	FindAll() (Rockets, error)
}

type RocketRepositoryWriter interface {
	Save(rocket *Rocket) error
	DeleteById(id RocketId) error
}
