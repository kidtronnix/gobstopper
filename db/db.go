// Package database is a package with convenience functions for working with SQL type dbs.
package database

import (
	"errors"
	"net/http"
	"strings"
	// Mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"github.com/jmoiron/sqlx"
)

// DatabaseKey is the key used to set the DB connection in the request context.
const DatabaseKey = "database"

type databaseName struct {
	Name string `db:"name"`
}

// ConnectToSQL is a function that will create a DB connection to your desired DB.
// Currently only the mysql driver is supported.
func ConnectToSQL(conn string) (*sqlx.DB, error) {
	connection := strings.Split(conn, "|")
	if len(connection) > 1 {
		return sqlx.Open(connection[0], connection[1])
	}

	return nil, errors.New("Bad connection string")
}

// GetConnection is used to get the *sqlx.DB connection from the context of a http request object.
// This only works if the connection has been set in the request context.
func GetConnection(r *http.Request) *sqlx.DB {
	db := context.Get(r, DatabaseKey).(*sqlx.DB)
	return db
}

// SetConnectionInRequestContext sets a *sqlx.DB connection in the context of the request using
// gorilla/context.
func SetConnectionInRequestContext(r *http.Request, db *sqlx.DB) {
	context.Set(r, DatabaseKey, db)
}

// GetDatabaseName will retrieve the database name for the currently connected database
func GetDatabaseName(db *sqlx.DB) (string, error) {
	dbName := databaseName{}
	err := db.Get(&dbName, "SELECT DATABASE() as name")
	return dbName.Name, err
}
