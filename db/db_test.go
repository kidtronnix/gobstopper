package database

import (
	"net/http"
	"os"
	"testing"

	"github.com/gorilla/context"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

var conn string

func init() {
	DB := os.Getenv("DB")
	switch DB {
	case "travis":
		conn = "root:@/golang?charset=utf8"
	default:
		conn = "root:root@/golang?charset=utf8"
	}
}

func TestSetDBInContext(t *testing.T) {
	assert := assert.New(t)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	assert.NoError(err)

	db, err := sqlx.Open("mysql", conn)
	assert.NoError(err)

	SetConnectionInRequestContext(req, db)

	assert.Equal(db, context.Get(req, DatabaseKey))
}

func TestGetFromContext(t *testing.T) {
	assert := assert.New(t)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	assert.NoError(err)

	db, err := sqlx.Open("mysql", conn)
	assert.NoError(err)

	context.Set(req, DatabaseKey, db)

	assert.Equal(db, GetConnection(req))
}

func TestConnectToSQL(t *testing.T) {
	assert := assert.New(t)

	db, err := ConnectToSQL("mysql|" + conn)
	assert.NoError(err)

	err = db.Ping()
	assert.NoError(err)

}

func TestConnectToSQLBadConnection(t *testing.T) {
	assert := assert.New(t)

	_, err := ConnectToSQL("imadoop")
	assert.Error(err)

}

func TestGetDatabaseName(t *testing.T) {
	assert := assert.New(t)

	db, err := sqlx.Open("mysql", conn)
	assert.NoError(err)

	name, err := GetDatabaseName(db)
	assert.NoError(err)

	assert.Equal("golang", name)
}
