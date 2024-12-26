package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ryanpujo/melius/config"
)

var db *sql.DB

func GetDBConnection() *sql.DB {
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()
	cfg := config.Config()

	counter := 0
	var err error

	for db == nil {
		db, err = sql.Open("pgx", cfg.DSN)
		if err != nil {
			log.Printf("DB connection attempt %d failed: %v", counter+1, err)
		}

		if err := db.Ping(); err != nil {
			log.Printf("failed to ping DB %d failed: %v", counter+1, err)
		}

		if counter == 5 {
			log.Fatalf("failed to connect to database after %d attempts", counter)
		}
		<-ticker.C
	}
	return db
}
