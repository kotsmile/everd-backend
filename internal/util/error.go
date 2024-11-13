package util

import "fmt"

type AppError struct {
	err         error
	internalErr error
}

func NewAppError(err error) AppError {
	return AppError{err: err}
}

func (e *AppError) WithError(err error) *AppError {
	e.internalErr = err
	return e
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.err, e.internalErr)
}
