package middleware

import (
	"context"
	"lunar-backend-engineer-challenge/pkg/logger"
	"lunar-backend-engineer-challenge/pkg/utils"
	"net/http"
	"strconv"

	"github.com/google/jsonapi"
	"go.uber.org/zap"
)

const (
	badRequestCode  = "bad_request"
	badRequestTitle = "Bad Request"

	conflictCode  = "conflict"
	conflictTitle = "Conflict"

	unauthorizedCode  = "unauthorized"
	unauthorizedTitle = "Unauthorized"

	forbiddenCode  = "forbidden"
	forbiddenTitle = "Forbidden"

	internalCode  = "internal_server_error"
	internalTitle = "Internal Server Error"

	notFoundCode  = "not_found"
	notFoundTitle = "Not Found"

	tooManyRequestsCode  = "too_many_requests"
	tooManyRequestsTitle = "Too Many Requests"

	headerXRequestID = "X-Request-ID"
)

type JsonApiResponseMiddleware struct {
	logger logger.Logger
}

func NewJsonApiResponseMiddleware(logger logger.Logger) *JsonApiResponseMiddleware {
	return &JsonApiResponseMiddleware{logger: logger}
}

func (jar *JsonApiResponseMiddleware) WriteErrorResponse(w http.ResponseWriter, errors []*jsonapi.ErrorObject, httpStatus int, previousError error) {
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(httpStatus)
	jar.logError(previousError, w.Header().Get(headerXRequestID), httpStatus)
	if err := jsonapi.MarshalErrors(w, errors); err != nil {
		jar.logger.Error("unexpected error marshalling json api response error", zap.Error(err), zap.String("correlation_id", w.Header().Get(headerXRequestID)))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (jar *JsonApiResponseMiddleware) WriteResponse(w http.ResponseWriter, payload interface{}, httpStatus int) {
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(httpStatus)

	if payload == nil {
		return
	}

	if err := jsonapi.MarshalPayload(w, payload); err != nil {
		jar.logger.Error("unexpected error marshalling json api response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (jar *JsonApiResponseMiddleware) logError(err error, correlationId string, statusCode int) {
	if err == nil {
		return
	}

	if statusCode >= http.StatusInternalServerError {
		logger.LogErrors(logger.Error, err, jar.logger, map[string]interface{}{"correlation_id": correlationId})
		return
	}
	logger.LogErrors(logger.Warning, err, jar.logger, map[string]interface{}{"correlation_id": correlationId})
}

func (jar *JsonApiResponseMiddleware) logErrorWithContext(ctx context.Context, err error, statusCode int) {
	if err == nil {
		return
	}

	if statusCode >= http.StatusInternalServerError {
		logger.LogErrorsWithContext(ctx, logger.Error, err, jar.logger, map[string]interface{}{})
		return
	}
	logger.LogErrorsWithContext(ctx, logger.Warning, err, jar.logger, map[string]interface{}{})
}

func BadRequestJsonApiHttpResponse(detail string) []*jsonapi.ErrorObject {
	return []*jsonapi.ErrorObject{{
		ID:     utils.NewUlid().String(),
		Code:   badRequestCode,
		Title:  badRequestTitle,
		Detail: detail,
		Status: strconv.Itoa(http.StatusBadRequest),
	}}
}

func ConflictJsonApiHttpResponse(detail string) []*jsonapi.ErrorObject {
	return []*jsonapi.ErrorObject{{
		ID:     utils.NewUlid().String(),
		Code:   conflictCode,
		Title:  conflictTitle,
		Detail: detail,
		Status: strconv.Itoa(http.StatusConflict),
	}}
}

func UnauthorizedRequestJsonApiHttpResponse(detail string) []*jsonapi.ErrorObject {
	return []*jsonapi.ErrorObject{{
		ID:     utils.NewUlid().String(),
		Code:   unauthorizedCode,
		Title:  unauthorizedTitle,
		Detail: detail,
		Status: strconv.Itoa(http.StatusUnauthorized),
	}}
}

func ForbiddenRequestJsonApiHttpResponse(detail string) []*jsonapi.ErrorObject {
	return []*jsonapi.ErrorObject{{
		ID:     utils.NewUlid().String(),
		Code:   forbiddenCode,
		Title:  forbiddenTitle,
		Detail: detail,
		Status: strconv.Itoa(http.StatusForbidden),
	}}
}

func InternalServerErrorJsonApiHttpResponse() []*jsonapi.ErrorObject {
	return []*jsonapi.ErrorObject{{
		ID:     utils.NewUlid().String(),
		Code:   internalCode,
		Title:  internalTitle,
		Detail: internalTitle,
		Status: strconv.Itoa(http.StatusInternalServerError),
	}}
}

func NotFoundRequestJsonApiHttpResponse(detail string) []*jsonapi.ErrorObject {
	return []*jsonapi.ErrorObject{{
		ID:     utils.NewUlid().String(),
		Code:   notFoundCode,
		Title:  notFoundTitle,
		Detail: detail,
		Status: strconv.Itoa(http.StatusNotFound),
	}}
}

func TooManyRequestsJsonApiHttpResponse(detail string) []*jsonapi.ErrorObject {
	return []*jsonapi.ErrorObject{{
		ID:     utils.NewUlid().String(),
		Code:   tooManyRequestsCode,
		Title:  tooManyRequestsTitle,
		Detail: detail,
		Status: strconv.Itoa(http.StatusTooManyRequests),
	}}
}
