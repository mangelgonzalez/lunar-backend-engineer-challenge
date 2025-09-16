package logger

import (
	"context"
	"fmt"
	"lunar-backend-engineer-challenge/pkg/domain"

	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
	"go.uber.org/zap"
)

const (
	Warning  = 1
	Error    = 2
	Previous = "previous_error"
)

func LogErrors(severity int, err error, logger Logger, items map[string]interface{}) {
	LogErrorsWithContext(context.Background(), severity, err, logger, items)
}

func LogErrorsWithContext(ctx context.Context, severity int, err error, logger Logger, items map[string]interface{}) {

	fields := append(extraItemsToFields(items), errorToContext(err)...)
	if severity == Warning {
		logger.Warn(err.Error(), fields...)
	} else {
		logger.Error(err.Error(), fields...)
	}
}

func extraItemsToFields(extraItems map[string]interface{}) []zap.Field {
	fields := make([]zap.Field, 0)
	for name, value := range extraItems {
		fields = append(fields, zap.Any(name, value))
	}

	return fields
}

func Log(err error, logger Logger) {
	LogWithContext(context.Background(), err, logger)
}

func LogWithContext(ctx context.Context, err error, logger Logger) {
	items := make(map[string]interface{})
	LogWithItemsContext(ctx, err, logger, items)
}

func LogWithItemsContext(ctx context.Context, err error, logger Logger, extraItems map[string]interface{}) {
	severity := Error

	if baseError, ok := errors.Cause(err).(domain.BaseError); ok {
		severity = baseError.Severity()
	}

	LogErrorsWithContext(ctx, severity, err, logger, extraItems)
}

func LogWithItems(err error, logger Logger, extraItems map[string]interface{}) {
	severity := Error

	if baseError, ok := errors.Cause(err).(domain.BaseError); ok {
		severity = baseError.Severity()
	}

	LogErrors(severity, err, logger, extraItems)
}

func errorToContext(err error) []zap.Field {
	items := errorToItems(err)
	return extraItemsToFields(items)
}

func formattedError(err error) map[string]interface{} {
	return funk.Union(extraItems(err), map[string]interface{}{
		"error": err.Error(),
	}).(map[string]interface{})
}

func extraItems(err error) map[string]interface{} {
	extraItems := make(map[string]interface{})
	if baseError, ok := errors.Cause(err).(domain.BaseError); ok {
		extraItems = baseError.ExtraItems()
	}
	return extraItems
}

func errorToItems(err error) map[string]interface{} {
	return funk.Union(formattedError(err), previousErrorsInfo(err)).(map[string]interface{})
}

func previousErrorsInfo(err error) map[string]interface{} {
	level := 0
	aggFunc := func(acc map[string]interface{}, erro error) map[string]interface{} {

		var newKey string
		if level > 0 {
			newKey = fmt.Sprintf("%s_%d", Previous, level)
		} else {
			newKey = Previous
		}
		level++

		return funk.Union(acc, map[string]interface{}{
			newKey: formattedError(erro),
		}).(map[string]interface{})

	}

	return funk.Reduce(allPreviousErrors(err), aggFunc, map[string]interface{}{}).(map[string]interface{})

}

func allPreviousErrors(err error) []error {
	var allErrors []error

	for {

		err = errors.Cause(err)
		baseError, ok := err.(domain.BaseError)
		if !ok {
			break
		}
		previous := baseError.Previous()
		if previous == nil {
			break
		}
		allErrors = append(allErrors, previous)
		err = previous
	}

	return allErrors
}

func ExtraItems(error error, fields []zap.Field) []zap.Field {
	error = errors.Cause(error)
	if baseError, ok := error.(domain.BaseError); ok {
		if previous := baseError.Previous(); previous != nil {
			fields = append(fields, ExtraItems(previous, fields)...)
		}
		for key, value := range baseError.ExtraItems() {
			fields = append(fields, zap.Reflect(key, value))
		}
	}

	fields = append(fields, zap.Error(error))

	return fields
}
