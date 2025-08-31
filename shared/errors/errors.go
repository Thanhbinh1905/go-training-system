package errors

import "errors"

var (
	ErrBadRequest   = errors.New("bad request")
	ErrUnauthorized = errors.New("unauthorized")   // 401
	ErrForbidden    = errors.New("forbidden")      // 403
	ErrNotFound     = errors.New("not found")      // 404
	ErrConflict     = errors.New("conflict")       // 409
	ErrInvalidInput = errors.New("invalid input")  // 400
	ErrInternal     = errors.New("internal error") // 500
)
