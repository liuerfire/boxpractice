package services

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound = errors.New("resource not found")
)

// ErrCode is one of a server-defined set of error codes.
type ErrCode string

const (
	ErrBadArgument      ErrCode = "BadArgument"
	ErrResourceNotFound ErrCode = "ResourceNotFound"
	ErrAlreadyExists    ErrCode = "ResourceAlreadyExists"
	ErrPermissionDenied ErrCode = "PermissionDenied"
	ErrInternalError    ErrCode = "InternalError"
)

type ServiceError struct {
	Code ErrCode
	Msg  string
}

func (s *ServiceError) Error() string {
	return string(s.Code)
}

// StatusCode indicates which HTTP response status code the errCode belongs to.
func (e ErrCode) StatusCode() int {
	switch e {
	case ErrBadArgument:
		return http.StatusBadRequest
	case ErrResourceNotFound:
		return http.StatusNotFound
	case ErrPermissionDenied:
		return http.StatusForbidden
	case ErrAlreadyExists:
		return http.StatusConflict
	}
	return http.StatusInternalServerError
}
