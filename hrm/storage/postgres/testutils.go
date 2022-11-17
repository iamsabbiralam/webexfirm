package postgres

import (
	"fmt"
	"log"
	"regexp"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
)

// replaceDBName replaces the dbname option in connection string with given db name in parameter.
func replaceDBName(connStr, dbName string) string {
	r := regexp.MustCompile(`dbname=([^\s]+)\s`)
	return r.ReplaceAllString(connStr, fmt.Sprintf("dbname=%s ", dbName))
}

// MustNewDevelopmentDB creates a new isolated database for the use of a package test
// The checking of dbconn is expected to be done in the package test using this
func MustNewDevelopmentDB(ddlConnStr, migrationDir string) (*sqlx.DB, func()) {
	const driver = "postgres"
	dbName := uuid.New().String()
	ddlDB := sqlx.MustConnect(driver, ddlConnStr)
	ddlDB.MustExec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbName))
	fmt.Println("Test Database name : ", dbName)
	if err := ddlDB.Close(); err != nil {
		panic(err)
	}
	connStr := replaceDBName(ddlConnStr, dbName)
	db := sqlx.MustConnect(driver, connStr)

	if err := goose.Run("up", db.DB, migrationDir); err != nil {
		panic(err)
	}

	tearDownFn := func() {
		if err := db.Close(); err != nil {
			log.Fatalf("failed to close database connection: %s", err.Error())
		}
		ddlDB, err := sqlx.Connect(driver, ddlConnStr)
		if err != nil {
			log.Fatalf("failed to connect database: %s", err.Error())
		}

		if _, err = ddlDB.Exec(fmt.Sprintf(`DROP DATABASE "%s"`, dbName)); err != nil {
			log.Fatalf("failed to drop database: %s", err.Error())
		}

		if err = ddlDB.Close(); err != nil {
			log.Fatalf("failed to close DDL database connection: %s", err.Error())
		}
	}
	return db, tearDownFn
}
