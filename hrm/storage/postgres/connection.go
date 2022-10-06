package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

const driver = "postgres"

// NewDBStringFromConfig build database connection string from config file.
func NewDBStringFromConfig(config *viper.Viper) (string, error) {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.GetString("database.user"),
		config.GetString("database.password"),
		config.GetString("database.host"),
		config.GetString("database.port"),
		config.GetString("database.dbname"),
		config.GetString("database.sslMode"),
	), nil
}

// Open opens a connection to database with given connection string.
func Open(config *viper.Viper) (*sql.DB, error) {
	dbString, err := NewDBStringFromConfig(config)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(driver, dbString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Open opens a connection to database with given connection string, using sqlx opener.
func Openx(config *viper.Viper) (*sqlx.DB, error) {
	dbString, err := NewDBStringFromConfig(config)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open(driver, dbString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Connectx opens a connection to database with given connection string using sqlx opener
// and verify the connection with a ping.
func Connectx(config *viper.Viper) (*sqlx.DB, error) {
	dbString, err := NewDBStringFromConfig(config)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Connect(driver, dbString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
