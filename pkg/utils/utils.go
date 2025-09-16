package utils

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
)

const UlidRegex = "[0123456789ABCDEFGHJKMNPQRSTVWXYZ]{26}"

type Ulid string

func (u Ulid) String() string {
	return string(u)
}

func NewUlid() Ulid {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return Ulid(ulid.MustNew(ulid.Timestamp(t), entropy).String())
}

func UlidFromString(rawUlid string) (Ulid, error) {
	var ulid Ulid
	if err := guardULID(rawUlid); err != nil {
		return ulid, err
	}

	return Ulid(rawUlid), nil
}

func guardULID(rawUlid string) error {
	_, err := ulid.Parse(rawUlid)
	return err
}

type Uuid string

func NewUuid() Uuid {
	return Uuid(uuid.New().String())
}

func (u Uuid) String() string {
	return string(u)
}

func GetMapValue(key string, myMap map[string]interface{}) interface{} {
	if x, found := myMap[key]; found {
		return x
	}

	return nil
}

func TimeToMilliseconds(someTime time.Time) int64 {
	return someTime.UnixNano() / int64(time.Millisecond)
}

func TimedBackoff(ctx context.Context, maxTryDuration time.Duration, fn func(ctx context.Context) error) error {
	var err error
	for {
		select {
		case <-ctx.Done():
			if err == nil {
				return errors.New("timed backoff failed without any iteration")
			}
			return err
		default:
			tryCtx, tryCn := context.WithTimeout(ctx, maxTryDuration)
			err = fn(tryCtx)
			tryCn()
			if err == nil {
				return nil
			}
		}
	}
}

func RetryFunc(fn func() (interface{}, error), timesToRetry int) (interface{}, error) {
	var (
		err    error
		result interface{}
	)

	for i := 1; i <= timesToRetry; i++ {
		result, err = fn()
		if err == nil {
			return result, nil
		}
	}

	return nil, err
}

func MapInterfaceInterfaceToStringInterface(payload map[interface{}]interface{}) map[string]interface{} {
	newPayload := make(map[string]interface{}, 0)
	for s, i := range payload {
		newPayload[fmt.Sprint(s)] = i
	}

	return newPayload
}

func GetInMapValueOrDefault(keyToFind []string, myMap map[string]interface{}, defaultValue interface{}) interface{} {
	current := myMap
	for i := 0; i < len(keyToFind); i++ {
		if _, found := current[keyToFind[i]]; !found {
			return defaultValue
		}

		if i+1 == len(keyToFind) {
			return current[keyToFind[i]]
		}

		current = current[keyToFind[i]].(map[string]interface{})
	}

	return current
}

func MillisToTime(millis int64) time.Time {
	return time.Unix(0, millis*int64(time.Millisecond))
}

func Pointer[T any](value T) *T {
	return &value
}
