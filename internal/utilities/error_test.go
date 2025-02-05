package utilities_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ryanpujo/melius/internal/utilities"
	"github.com/stretchr/testify/require"
)

func TestHandleError(t *testing.T) {
	tableTest := map[string]struct {
		err    error
		assert func(t *testing.T, err error)
	}{
		"sql no rows error": {
			err: sql.ErrNoRows,
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				var appErr *utilities.AppError
				require.ErrorAs(t, err, &appErr)
				require.Equal(t, utilities.DataNotFound, appErr.Code)
				require.Equal(t,
					"We're sorry, but we couldn't find any data with that information. Please check your input and try again.",
					appErr.Message,
				)
			},
		},
		"pg error": {
			err: &pgconn.PgError{
				Code: utilities.PGCheckViolation,
			},
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				var appErr *utilities.AppError
				require.ErrorAs(t, err, &appErr)
				require.Equal(t, utilities.CheckViolation, appErr.Code)
				require.Equal(t,
					"The provided data does not meet the required criteria.",
					appErr.Message,
				)
			},
		},
		"validation error": {
			err: validator.ValidationErrors{},
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				var appErr *utilities.AppError
				require.ErrorAs(t, err, &appErr)
				require.Equal(t, utilities.ValidationError, appErr.Code)
				require.Equal(t,
					"There was a problem with your request. Please double-check your input and try again.",
					appErr.Message,
				)
			},
		},
		"unknown error": {
			err: errors.New("failed"),
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				var appErr *utilities.AppError
				require.ErrorAs(t, err, &appErr)
				require.Equal(t, utilities.InternalServerError, appErr.Code)
				require.Equal(t,
					"An internal server error occurred. Please try again later.",
					appErr.Message,
				)
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			err := utilities.HandleError(v.err)

			v.assert(t, err)
		})
	}
}
