package sync

import "lunar-backend-engineer-challenge/pkg/domain"

const (
	errorLockingKeyMessage   = "Something happened while acquiring processes"
	errorReleasingKeyMessage = "Something happens while releasing processes"
)

type ErrorMutex struct {
	domain.CriticalError
	message    string
	previous   error
	extraItems map[string]interface{}
}

func (i ErrorMutex) ExtraItems() map[string]interface{} {
	return i.extraItems
}

func (i ErrorMutex) Previous() error {
	return i.previous
}
func (i ErrorMutex) Unwrap() error {
	return i.previous
}

func (i ErrorMutex) Error() string {
	return i.message
}

func NewErrorLockMutexKey(identifier string, previous error) *ErrorMutex {
	return &ErrorMutex{message: errorLockingKeyMessage, extraItems: map[string]interface{}{"identifier": identifier}, previous: previous}
}

func NewErrorReleaseLockMutexKey(identifier string, previous error) *ErrorMutex {
	return &ErrorMutex{message: errorReleasingKeyMessage, extraItems: map[string]interface{}{"identifier": identifier}, previous: previous}
}
