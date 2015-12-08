package database

import (
	"net/http"
	"testing"

	"github.com/gorilla/context"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSetDBInContext(t *testing.T) {
	assert := assert.New(t)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	assert.NoError(err)

	db, err := sqlx.Open("mysql", "root:root@/golang?charset=utf8&parseTime=True&loc=Local")
	assert.NoError(err)

	SetConnectionInRequestContext(req, db)

	assert.Equal(db, context.Get(req, DatabaseKey))
}

func TestGetFromContext(t *testing.T) {
	assert := assert.New(t)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	assert.NoError(err)

	db, err := sqlx.Open("mysql", "root:root@/golang?charset=utf8&parseTime=True&loc=Local")
	assert.NoError(err)

	context.Set(req, DatabaseKey, db)

	assert.Equal(db, GetConnection(req))
}

func TestConnectToSQL(t *testing.T) {
	assert := assert.New(t)

	db, err := ConnectToSQL("mysql|root:root@/golang?charset=utf8&parseTime=True&loc=Local")
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

	db, err := sqlx.Open("mysql", "root:root@/golang?charset=utf8&parseTime=True&loc=Local")
	assert.NoError(err)

	name, err := GetDatabaseName(db)
	assert.NoError(err)

	assert.Equal("golang", name)
}
