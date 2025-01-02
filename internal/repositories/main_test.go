package repositories_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryanpujo/melius/internal/repositories"
)

var (
	db             *sql.DB
	mock           sqlmock.Sqlmock
	credentialRepo *repositories.CredentialRepo
	addressRepo    repositories.AddressRepo
)

func TestMain(m *testing.M) {
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	credentialRepo = repositories.NewCredentialRepo(db)
	addressRepo = repositories.NewAddressRepo(db)

	os.Exit(m.Run())
}
