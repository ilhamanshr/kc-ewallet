package migrations

import (
	"database/sql"
	"fmt"
	log_color "kc-ewallet/internals/helpers/color"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func MigrateAll(db *sql.DB, dbName string) error {
	log_color.PrintYellow("Migrating pending migrations...")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("error on initiating postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://./migrations", dbName, driver)
	if err != nil {
		return fmt.Errorf("error on NewWithDatabaseInstance(): %v", err)
	}

	err = m.Up()
	if err != nil && err == migrate.ErrNoChange {
		slog.Info("There is no pending migration")
		return nil
	}

	if err != nil && err != migrate.ErrNilVersion && err != migrate.ErrNoChange {
		return fmt.Errorf("error when running migrate up: %v", err)
	}

	log_color.PrintGreen("All migration has been migrated successfully!")

	return nil
}
