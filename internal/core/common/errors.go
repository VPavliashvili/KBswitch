package common

import (
	"encoding/json"
	"errors"
	"net/http"
)

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	json, _ := json.Marshal(e)
	return string(json[:])
}

var (
	ErrBadRequest     = errors.New("bad request!")
	ErrNotFound       = errors.New("not found!")
	ErrInternalServer = errors.New("internal server error!")
)

type AppError struct {
	Errtype error
	Reason  error
}

func (e AppError) Error() string {
	return errors.Join(e.Errtype, e.Reason).Error()
}

func NewError(errtype error, reason string) AppError {
	return AppError{
		Errtype: errtype,
		Reason:  errors.New(reason),
	}
}

func ToAPIErr(err AppError) APIError {
	var status int
	switch err.Errtype {
	case ErrBadRequest:
		status = http.StatusBadRequest
	case ErrNotFound:
		status = http.StatusNotFound
	case ErrInternalServer:
		status = http.StatusInternalServerError
	}

	return APIError{Status: status, Message: err.Reason.Error()}
}
