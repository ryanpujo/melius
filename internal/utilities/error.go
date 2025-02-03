package utilities

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	_ = iota
	DataNotFound
	NotNullViolation
	DuplicateData
	CheckViolation
	InternalServerError
)

type ErrorCode int

func (err ErrorCode) String() string {
	return [...]string{
		"", // index 0 (unused due to iota skip)
		"Data not found",
		"Not null violation",
		"Duplicate data",
		"Check constraint violation",
		"Internal Server Error",
	}[err]
}

func (err ErrorCode) HTTPStatus() int {
	switch err {
	case DataNotFound:
		return http.StatusNotFound
	case NotNullViolation, CheckViolation:
		return http.StatusBadRequest
	case DuplicateData:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Err     error     `json:"-"`
}

func (ae *AppError) Error() string {
	if ae.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", ae.Code.String(), ae.Message, ae.Err)
	}
	return fmt.Sprintf("[%s] %s", ae.Code.String(), ae.Message)
}

func (ae *AppError) Unwrap() error {
	return ae.Err
}

func NewAppError(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func FromPGError(err *pgconn.PgError) *AppError {
	switch err.Code {
	case "23502":
		// NOT_NULL_VIOLATION
		return NewAppError(NotNullViolation, "A required field is missing. Please check your input.", err)
	case "23505":
		// duplicate data
		return NewAppError(DuplicateData, "This data already exists. Please provide a unique value.", err)
	case "23514":
		// check violation
		return NewAppError(CheckViolation, "The provided data does not meet the required criteria.", err)
	default:
		return NewAppError(InternalServerError, "An internal error occurred. Please try again later.", err)
	}
}
