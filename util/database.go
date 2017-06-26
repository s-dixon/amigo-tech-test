package util

import "database/sql"

type DatabaseConnector interface {
	Open(connectionString string) (*sql.DB, error)
}

type PostgresDatabaseConnector struct {}

func (PostgresDatabaseConnector) Open(connectionString string) (*sql.DB, error) {
	return sql.Open("postgres", connectionString)
}
