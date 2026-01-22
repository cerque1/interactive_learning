package usecase

import (
	"errors"
	"fmt"
)

var NotFoundErr = errors.New("not found")

type NotFoundRecordErr struct {
	Err error
}

func NewNotFoundErr(err error) *NotFoundRecordErr {
	return &NotFoundRecordErr{Err: err}
}

func (nf *NotFoundRecordErr) Error() string {
	return nf.Err.Error()
}

func (nf *NotFoundRecordErr) Unwrap() error {
	return NotFoundErr
}

var InternalErr = errors.New("internal server error")

type InternalError struct {
	Err error
}

func NewInternalError(err error) *InternalError {
	return &InternalError{Err: err}
}

func (idbe *InternalError) Error() string {
	return idbe.Err.Error()
}

func (idbe *InternalError) Unwrap() error {
	return InternalErr
}

var UnauthorizedErr = errors.New("Unauthorized user")

type UnauthorizedError struct {
	Err error
}

func NewUnauthorizedError(err error) *UnauthorizedError {
	return &UnauthorizedError{Err: err}
}

func (ue *UnauthorizedError) Error() string {
	return ue.Err.Error()
}

func (ue *UnauthorizedError) Unwrap() error {
	return UnauthorizedErr
}

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

var ChangeTypeErr = errors.New("change type error")

type ChangeTypeError struct {
	Object string
	Err    error
}

func NewChangeTypeError(object string, err error) *ChangeTypeError {
	return &ChangeTypeError{Object: object, Err: err}
}

func (cte *ChangeTypeError) Error() string {
	return fmt.Errorf("change type error object: %s, error: %w", cte.Object, cte.Err).Error()
}

func (cte *ChangeTypeError) Unwrap() error {
	return ChangeTypeErr
}

var AlreadyExistsErr = errors.New("object is already exists")

type AlreadyExistsError struct {
	Object   string
	ObjectId int
}

func NewAlreadyExistsError(object string, objectId int) *AlreadyExistsError {
	return &AlreadyExistsError{Object: object, ObjectId: objectId}
}

func (ae *AlreadyExistsError) Error() string {
	return fmt.Sprintf("error: %s is already exists with id %d", ae.Object, ae.ObjectId)
}

func (ae *AlreadyExistsError) Unwrap() error {
	return AlreadyExistsErr
}
