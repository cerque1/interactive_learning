package repo

import (
	"errors"
	"fmt"
)

// Отсутсвие записи для удаления
var NoSuchRecordToDelete = errors.New("no such record to delete")
var NoSuchRecordToUpdate = errors.New("no such record to update")
var NoSuchRecordToSelect = errors.New("no such record to select")

// Ошибка добавления записи
var InsertRecordError = errors.New("insert record error")

// При ошибке при выполнении sql запроса
var DBErr = errors.New("psql error")

type DBError struct {
	Table     string
	Operation string
	Err       error
}

func NewDBError(table string, operation string, err error) *DBError {
	return &DBError{Table: table, Operation: operation, Err: err}
}

func (dbe *DBError) Error() string {
	return fmt.Sprintf("psql error: table: %s, operation: %s, error: %s", dbe.Table, dbe.Operation, dbe.Err.Error())
}

func (dbe *DBError) Unwrap() error {
	return DBErr
}

var InvalidToken = errors.New("invalid token")
var ExpiredToken = errors.New("token is expired")
