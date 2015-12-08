package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServiceWithDB(t *testing.T) {
	assert := assert.New(t)

	service, err := NewService(8000, "/v1/prefix", "mysql|root:root@/golang?charset=utf8&parseTime=True&loc=Local")
	assert.NoError(err)
	assert.NotNil(service)
	assert.IsType(&Service{}, service)
	assert.Equal(8000, service.port)
	assert.Equal("/v1/prefix", service.prefix)
	assert.NotNil(service.Router)
	assert.Empty(service.routes)
	assert.Empty(service.middleware)
	err = service.DB.Ping()
	assert.NoError(err)

	service.AddRoute("GET", "", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		fmt.Fprint(rw, "ok")
	})

	m := func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		rw.Header().Set("Foo", "bar")
		next(rw, r)
	}
	service.AddMiddleware(service.DBConnectionMiddleware)
	service.AddMiddleware(m)

	go service.Start()

	time.Sleep(250 * time.Millisecond)

	response, err := http.Get("http://localhost:8000/v1/prefix")

	assert.NoError(err)
	if err == nil {
		defer response.Body.Close()
		assert.Equal("bar", response.Header.Get("Foo"))
		assert.Equal(200, response.StatusCode)
		body, err := ioutil.ReadAll(response.Body)
		assert.NoError(err)
		assert.Equal("ok", string(body))
	}
}

func TestServiceWithoutDB(t *testing.T) {
	assert := assert.New(t)

	service, err := NewService(8000, "/v1/prefix", "")
	assert.NoError(err)
	assert.NotNil(service)
	assert.IsType(&Service{}, service)
	assert.Equal(8000, service.port)
	assert.Equal("/v1/prefix", service.prefix)
	assert.NotNil(service.Router)
	assert.Empty(service.routes)
	assert.Empty(service.middleware)
	assert.Nil(service.DB)
}

func TestServiceErrorsOnConnectionTypeFailure(t *testing.T) {
	assert := assert.New(t)

	service, err := NewService(8000, "/v1/prefix", "sadasd")
	assert.Error(err)
	assert.Nil(service)
}

func TestServiceErrorsOnConnectionFailure(t *testing.T) {
	assert := assert.New(t)

	service, err := NewService(8000, "/v1/prefix", "mysql|asdas")
	assert.Error(err)
	assert.Nil(service)
}
