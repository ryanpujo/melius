package repositories_test

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryanpujo/melius/internal/models"
	"github.com/ryanpujo/melius/internal/repositories"
	"github.com/stretchr/testify/require"
)

var (
	db             *sql.DB
	mock           sqlmock.Sqlmock
	credentialRepo *repositories.CredentialRepo
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
	credential := models.CredentialPayload{
		Email:    "ryanpujo@gmail.com",
		Username: "ryanpujo",
		Password: "okeoke",
	}
	user := models.UserPayload{
		FirstName:         "Ryan",
		LastName:          "Pujo",
		CredentialPayload: credential,
	}
	tableTest := map[string]struct {
		arrange func()
		assert  func(t *testing.T, id uint, err error)
	}{
		"success": {
			arrange: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO credentials").
					WithArgs(
						credential.Email,
						credential.Username,
						credential.Password,
						time.Now().Format(time.RFC3339),
						time.Now().Format(time.RFC3339)).
					WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow(credential.Username))

				mock.ExpectQuery("INSERT INTO users").
					WithArgs(
						user.FirstName,
						user.LastName,
						user.CredentialPayload.Username,
						time.Now().Format(time.RFC3339),
						time.Now().Format(time.RFC3339)).
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
						credential.Email,
						credential.Username,
						credential.Password,
						time.Now().Format(time.RFC3339),
						time.Now().Format(time.RFC3339)).
					WillReturnError(errors.New("failed"))

				mock.ExpectQuery("INSERT INTO users").
					WithArgs(
						user.FirstName,
						user.LastName,
						user.CredentialPayload.Username,
						time.Now().Format(time.RFC3339),
						time.Now().Format(time.RFC3339)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

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
						credential.Email,
						credential.Username,
						credential.Password,
						time.Now().Format(time.RFC3339),
						time.Now().Format(time.RFC3339)).
					WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow(credential.Username))

				mock.ExpectQuery("INSERT INTO users").
					WithArgs(
						user.FirstName,
						user.LastName,
						user.CredentialPayload.Username,
						sqlmock.AnyArg(),
						sqlmock.AnyArg()).
					WillReturnError(errors.New("user failed"))

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

	for key, v := range tableTest {
		t.Run(key, func(t *testing.T) {
			v.arrange()

			id, err := credentialRepo.Write(context.Background(), user)

			v.assert(t, id, err)
		})
	}
}
