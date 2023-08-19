package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/slavik22/imageAPI/api"
	db "github.com/slavik22/imageAPI/db/sqlc"
	"github.com/slavik22/imageAPI/util"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("fatal error config file: %w", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot open db connection ", err)
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal("Cannot create server ", err)
	}

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot start server ", err)
	}
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot create new migrate instance %e", err))
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(fmt.Errorf("failed to run migrate up %e", err))
	}

	log.Println("db migrated successfully")
}
