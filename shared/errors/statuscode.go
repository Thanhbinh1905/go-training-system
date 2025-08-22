package errors

import (
	"errors"
	"net/http"
)

func StatusCode(err error) int {
	switch {
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized // 401
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden // 403
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound // 404
	case errors.Is(err, ErrConflict):
		return http.StatusConflict // 409
	case errors.Is(err, ErrInvalidInput):
		return http.StatusBadRequest // 400
	default:
		return http.StatusInternalServerError // 500
	}
}
