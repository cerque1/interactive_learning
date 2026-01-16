package errors

import (
	"errors"
	"fmt"
)

var ErrNotAvailable = errors.New("error: object not available")

type NotAvailable struct {
	objectType string
	objectId   int
}

func NewNotAvailableError(objectType string, objectId int) *NotAvailable {
	return &NotAvailable{objectType: objectType, objectId: objectId}
}

func (na *NotAvailable) Error() string {
	return fmt.Sprintf("error: the object %s with id %d is unavailable", na.objectType, na.objectId)
}

func (na *NotAvailable) Unwrap() error {
	return ErrNotAvailable
}
