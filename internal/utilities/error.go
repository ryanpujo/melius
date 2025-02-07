package utilities

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
)

// PostgreSQL error codes as constants for better readability
const (
	PGNotNullViolation  = "23502"
	PGDuplicateData     = "23505"
	PGCheckViolation    = "23514"
)

const (
	_ ErrorCode = iota
	DataNotFound
	NotNullViolation
	DuplicateData
	CheckViolation
	InternalServerError
	DatabaseError
	ValidationError
)

type ErrorCode int

// errorCodeStrings now uses a map for safer lookups and easier maintenance
var errorCodeStrings = map[ErrorCode]string{
	DataNotFound:        "Data not found",
	NotNullViolation:    "Not null violation",
	DuplicateData:       "Duplicate data",
	CheckViolation:      "Check constraint violation",
	InternalServerError: "Internal Server Error",
	DatabaseError:       "Internal Database Error",
	ValidationError:     "Validation Error",
}

func (err ErrorCode) String() string {
	if s, ok := errorCodeStrings[err]; ok {
		return s
	}
	return "Unknown error"
}

func (err ErrorCode) HTTPStatus() int {
	switch err {
	case DataNotFound:
		return http.StatusNotFound
	case NotNullViolation, CheckViolation, ValidationError:
		return http.StatusBadRequest
	case DuplicateData:
		return http.StatusConflict
	case InternalServerError, DatabaseError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
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
	case PGNotNullViolation:
		return NewAppError(NotNullViolation, "A required field is missing. Please check your input.", err)
	case PGDuplicateData:
		return NewAppError(DuplicateData, "This data already exists. Please provide a unique value.", err)
	case PGCheckViolation:
		return NewAppError(CheckViolation, "The provided data does not meet the required criteria.", err)
	default:
		return NewAppError(DatabaseError, "An internal database error occurred. Please try again later.", err)
	}
}

func HandleError(err error) *AppError {
	var appErr *AppError
	var pgErr *pgconn.PgError
	var valErr validator.ValidationErrors

	switch {
	case errors.Is(err, sql.ErrNoRows):
		appErr = NewAppError(
			DataNotFound,
			"We're sorry, but we couldn't find any data with that information. Please check your input and try again.",
			err,
		)
	case errors.As(err, &pgErr):
		appErr = FromPGError(pgErr)
	case errors.As(err, &valErr):
		appErr = NewAppError(
			ValidationError,
			"There was a problem with your request. Please double-check your input and try again.", 
			err,
		)
	default:
		appErr = NewAppError(
			InternalServerError,
			"An internal server error occurred. Please try again later.",
			err,
		)
	}

	return appErr
}