package main

import (
	"database/sql"
	"janan_csv_service/pkg/helpers"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("Starting migration")
	dbUrl := helpers.SafeGetEnv("LOCAL_DATABASE_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Panic(err)
		return
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:cmd/database/migrations",
		"postgres", driver)
	if err != nil {
		log.Panic(err)
		return
	}
	err = m.Up()
	if err != nil {
		log.Panic(err)
		return
	}
	log.Println("Migration completed")
}
