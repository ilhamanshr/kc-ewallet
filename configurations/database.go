package configurations

import (
	"os"
)

type databaseConfiguration struct {
	user     string
	password string
	host     string
	port     string
	dbname   string
	param    string
}

//go:generate mockgen -destination=mocks/mock_database.go -source=database.go IDatabaseConfiguration
type IDatabaseConfiguration interface {
	GetUser() string
	GetPassword() string
	GetHost() string
	GetPort() string
	GetDBName() string
	GetDBParam() string
}

func NewDatabaseWriter() *databaseConfiguration {
	return &databaseConfiguration{
		user:     os.Getenv("DB_USERNAME"),
		password: os.Getenv("DB_PASSWORD"),
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
		dbname:   os.Getenv("DB_NAME"),
		param:    os.Getenv("DB_PARAMS"),
	}
}

func (dw *databaseConfiguration) GetUser() string {
	return dw.user
}

func (dw *databaseConfiguration) GetPassword() string {
	return dw.password
}

func (dw *databaseConfiguration) GetHost() string {
	return dw.host
}

func (dw *databaseConfiguration) GetPort() string {
	return dw.port
}

func (dw *databaseConfiguration) GetDBName() string {
	return dw.dbname
}

func (dw *databaseConfiguration) GetDBParam() string {
	return dw.param
}
