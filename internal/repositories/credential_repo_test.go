package repositories_test

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryanpujo/melius/internal/models"
	"github.com/ryanpujo/melius/internal/repositories"
	"github.com/stretchr/testify/require"
)

var (
	db                *sql.DB
	mock              sqlmock.Sqlmock
	credentialRepo    *repositories.CredentialRepo
	credentialPayload = models.CredentialPayload{
		Email:    "ryanpujo@gmail.com",
		Username: "ryanpujo",
		Password: "okeoke",
	}
	credential = &models.Credential{
		Email:    "ryanpujo@gmail.com",
		Username: "ryanpujo",
		Password: "okeoke",
	}
	userPayload = models.UserPayload{
		FirstName:         "Ryan",
		LastName:          "Pujo",
		CredentialPayload: credentialPayload,
	}
	user = models.User{
		FirstName:  "Ryan",
		LastName:   "Pujo",
		Credential: *credential,
	}
)

func TestMain(m *testing.M) {
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	credentialRepo = repositories.NewCredentialRepo(db)

	os.Exit(m.Run())
}

func TestShouldCreateUserAndCredential(t *testing.T) {
	tableTest := map[string]struct {
		arrange func()
		assert  func(t *testing.T, id uint, err error)
	}{
		"success": {
			arrange: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO credentials").
					WithArgs(
						credentialPayload.Email,
						credentialPayload.Username,
						credentialPayload.Password,
						sqlmock.AnyArg(),
						sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow(credentialPayload.Username))

				mock.ExpectQuery("INSERT INTO users").
					WithArgs(
						userPayload.FirstName,
						userPayload.LastName,
						userPayload.CredentialPayload.Username,
						sqlmock.AnyArg(),
						sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectCommit()
			},
			assert: func(t *testing.T, id uint, err error) {
				require.NoError(t, err)
				require.Equal(t, uint(1), id)
			},
		},
		"rollback on credential failure": {
			arrange: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO credentials").
					WithArgs(
						credentialPayload.Email,
						credentialPayload.Username,
						credentialPayload.Password,
						sqlmock.AnyArg(),
						sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"username"}).RowError(1, errors.New("failed")))

				mock.ExpectRollback()
			},
			assert: func(t *testing.T, id uint, err error) {
				require.Error(t, err)
				require.Zero(t, id)
			},
		},
		"rollback on user failure": {
			arrange: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO credentials").
					WithArgs(
						credentialPayload.Email,
						credentialPayload.Username,
						credentialPayload.Password,
						sqlmock.AnyArg(),
						sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow(credentialPayload.Username))

				mock.ExpectQuery("INSERT INTO users").
					WithArgs(
						userPayload.FirstName,
						userPayload.LastName,
						userPayload.CredentialPayload.Username,
						sqlmock.AnyArg(),
						sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("dgrg"))

				mock.ExpectRollback()
			},
			assert: func(t *testing.T, id uint, err error) {
				require.Error(t, err)
				require.Zero(t, id)
			},
		},
		"faile to start transaction": {
			arrange: func() {
				mock.ExpectBegin().WillReturnError(errors.New("failed"))
			},
			assert: func(t *testing.T, id uint, err error) {
				require.Error(t, err)
				require.Zero(t, id)
			},
		},
	}

	for _, v := range tableTest {
		v.arrange()

		id, err := credentialRepo.Write(context.Background(), userPayload)

		v.assert(t, id, err)
		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	}
}

func TestFindByUsername(t *testing.T) {
	tableTest := map[string]struct {
		arrange func()
		assert  func(t *testing.T, actual *models.User, err error)
	}{
		"success": {
			arrange: func() {
				row := sqlmock.NewRows([]string{"first_name", "last_name", "email", "username", "password"}).
					AddRow(
						user.FirstName, user.LastName, user.Credential.Email, user.Credential.Username,
						user.Credential.Password,
					)

				mock.ExpectQuery(`
					SELECT u.first_name, u.last_name, c.email, c.username, c.password
					FROM users u
					JOIN credentials c ON c.username = u.username
					WHERE u.username = \$1
				`).WithArgs(credentialPayload.Username).WillReturnRows(row)
			},
			assert: func(t *testing.T, actual *models.User, err error) {
				require.NoError(t, err)
				require.Equal(t, &user, actual)
			},
		},
		"scan failed": {
			arrange: func() {
				row := sqlmock.NewRows([]string{"first_name", "last_name", "email", "username", "password"}).
					RowError(1, errors.New("failed to scan"))

				mock.ExpectQuery(`
					SELECT u.first_name, u.last_name, c.email, c.username, c.password
					FROM users u
					JOIN credentials c ON c.username = u.username
					WHERE u.username = \$1
				`).WithArgs(credentialPayload.Username).WillReturnRows(row)
			},
			assert: func(t *testing.T, actual *models.User, err error) {
				require.Error(t, err)
				require.Nil(t, actual)
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			cred, err := credentialRepo.FindByUsername(context.Background(), credentialPayload.Username)

			v.assert(t, cred, err)
			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Unmet expectations: %v", err)
			}
		})
	}
}
