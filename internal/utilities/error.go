package utilities

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
)

// PostgreSQL error codes as constants for better readability.
const (
	PGNotNullViolation = "23502" // Not-null constraint violation.
	PGDuplicateData    = "23505" // Unique constraint violation (duplicate data).
	PGCheckViolation   = "23514" // Check constraint violation.
)

// ErrorCode defines a custom type for error codes used in AppError.
type ErrorCode int

// Custom error codes used in the application.
const (
	_ ErrorCode = iota // Skip zero value.
	DataNotFound
	NotNullViolation
	DuplicateData
	CheckViolation
	InternalServerError
	DatabaseError
	ValidationError
)

// errorCodeStrings maps our custom error codes to their string representations.
// Using a map allows for safer lookups and easier maintenance.
var errorCodeStrings = map[ErrorCode]string{
	DataNotFound:        "Data not found",
	NotNullViolation:    "Not null violation",
	DuplicateData:       "Duplicate data",
	CheckViolation:      "Check constraint violation",
	InternalServerError: "Internal Server Error",
	DatabaseError:       "Internal Database Error",
	ValidationError:     "Validation Error",
}

// String returns the string representation of an ErrorCode.
func (err ErrorCode) String() string {
	if s, ok := errorCodeStrings[err]; ok {
		return s
	}
	return "Unknown error"
}

// HTTPStatus maps an ErrorCode to an appropriate HTTP status code.
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

// AppError is a custom error type that holds an error code, a human-readable message,
// and an underlying error that provides additional context.
type AppError struct {
	Code    ErrorCode // Custom error code.
	Message string    // Human-readable error message.
	Err     error     // Underlying error (if any).
}

// Error implements the error interface for AppError.
// It returns a formatted error string including the error code, message, and underlying error.
func (ae *AppError) Error() string {
	if ae.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", ae.Code.String(), ae.Message, ae.Err)
	}
	return fmt.Sprintf("[%s] %s", ae.Code.String(), ae.Message)
}

// Unwrap returns the underlying error, allowing errors.Is and errors.As to function properly.
func (ae *AppError) Unwrap() error {
	return ae.Err
}

// NewAppError creates a new AppError given a code, message, and an optional underlying error.
func NewAppError(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// FromPGError converts a PostgreSQL error (of type *pgconn.PgError) into an AppError.
// It uses specific SQLSTATE codes to determine the error type.
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

// HandleError inspects an error and returns an appropriate AppError.
// It handles common error cases such as no data found, PostgreSQL errors, and validation errors.
func HandleError(err error) *AppError {
	var appErr *AppError
	var pgErr *pgconn.PgError
	var valErr validator.ValidationErrors

	switch {
	// Handle the case where no rows are returned.
	case errors.Is(err, sql.ErrNoRows):
		appErr = NewAppError(
			DataNotFound,
			"We're sorry, but we couldn't find any data with that information. Please check your input and try again.",
			err,
		)
	// Check if the error is a PostgreSQL error and convert it accordingly.
	case errors.As(err, &pgErr):
		appErr = FromPGError(pgErr)
	// Check if the error is a validation error from the go-playground/validator package.
	case errors.As(err, &valErr):
		appErr = NewAppError(
			ValidationError,
			"There was a problem with your request. Please double-check your input and try again.",
			err,
		)
	// Fallback for all other errors.
	default:
		appErr = NewAppError(
			InternalServerError,
			"An internal server error occurred. Please try again later.",
			err,
		)
	}

	return appErr
}