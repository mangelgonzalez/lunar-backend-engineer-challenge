package domain

const (
	criticalError = 2
	domainError   = 1
)

type BaseError interface {
	Error() string
	ExtraItems() map[string]interface{}
	Severity() int
	Previous() error
}

type CriticalError struct {
	previous error
}

func (d *CriticalError) Severity() int {
	return criticalError
}

func (d *CriticalError) Previous() error {
	return d.previous
}

type DomainError struct {
	previous error
}

func (d *DomainError) Severity() int {
	return domainError
}

func (d *DomainError) Previous() error {
	return d.previous
}
