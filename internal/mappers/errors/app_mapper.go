package errors_mapper

import (
	"errors"
	"interactive_learning/internal/usecase"
	"net/http"
)

type ApplicationErrorsMapper struct{}

func NewApplicationErrorsMapper() *ApplicationErrorsMapper {
	return &ApplicationErrorsMapper{}
}

func (aem *ApplicationErrorsMapper) ApplicationErrorToHttp(err error) (int, map[string]string) {
	var answerStatus int
	switch {
	case errors.Is(err, usecase.NotFoundErr):
		answerStatus = http.StatusNotFound
	case errors.Is(err, usecase.UnauthorizedErr):
		answerStatus = http.StatusUnauthorized
	case errors.Is(err, usecase.ErrNotAvailable):
		answerStatus = http.StatusNotAcceptable
	case errors.Is(err, usecase.ChangeTypeErr),
		errors.Is(err, usecase.AlreadyExistsErr):
		answerStatus = http.StatusBadRequest
	default:
		answerStatus = http.StatusInternalServerError
	}

	return answerStatus, map[string]string{
		"message": err.Error(),
	}
}
