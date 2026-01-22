package errors_mapper

import (
	"errors"
	"interactive_learning/internal/repo"
	"interactive_learning/internal/usecase"
)

type DomainsErrorsMapper struct{}

func NewDomainErrorsMapper() *DomainsErrorsMapper {
	return &DomainsErrorsMapper{}
}

func (em *DomainsErrorsMapper) DBErrorToApp(err error) error {
	switch {
	case errors.Is(err, repo.NoSuchRecordToSelect),
		errors.Is(err, repo.NoSuchRecordToUpdate),
		errors.Is(err, repo.NoSuchRecordToDelete):
		return usecase.NewNotFoundErr(err)
	case errors.Is(err, repo.ExpiredToken),
		errors.Is(err, repo.InvalidToken):
		return usecase.NewUnauthorizedError(err)
	default:
		return usecase.NewInternalError(err)
	}
}
