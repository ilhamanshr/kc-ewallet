package database

import (
	"database/sql"
	"fmt"
	"kc-ewallet/configurations"
	"strings"
	"time"

	"github.com/XSAM/otelsql"
	_ "github.com/lib/pq" // PostgreSQL driver
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

type postgresWriter struct {
	db *sql.DB
}

type IPostgresWriter interface {
	GetDB() *sql.DB
	Close() error
}

func NewPostgresWriter(configDB configurations.IDatabaseConfiguration) (*postgresWriter, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		strings.TrimSpace(configDB.GetUser()),
		strings.TrimSpace(configDB.GetPassword()),
		strings.TrimSpace(configDB.GetHost()),
		strings.TrimSpace(configDB.GetPort()),
		strings.TrimSpace(configDB.GetDBName()),
		strings.TrimSpace(configDB.GetDBParam()))

	db, err := otelsql.Open("postgres", dsn, otelsql.WithAttributes(semconv.DBSystemNamePostgreSQL))
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(10 * time.Second)
	db.SetMaxOpenConns(90)
	if err != nil {
		// log.Fatal(err) // TODO: handle this properly!
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return &postgresWriter{db: db}, nil
}

func (mw *postgresWriter) GetDB() *sql.DB {
	return mw.db
}

func (mw *postgresWriter) Close() error {
	return mw.db.Close()
}
