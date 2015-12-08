package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitialisesFlags(t *testing.T) {
	assert := assert.New(t)

	arguments := []string{"-port", "8008", "-prefix", "/v2/stats", "-connection", "123:123@blah.com"}
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	initFlags(f, arguments)

	assert.Equal(8008, port)
	assert.Equal("/v2/stats", prefix)
	assert.Equal("123:123@blah.com", connection)

}

func TestGobstopperFunctional(t *testing.T) {
	assert := assert.New(t)

	arguments := []string{"-port", "8008", "-prefix", "/v1", "-connection", "mysql|root:root@/golang?charset=utf8"}
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	initFlags(f, arguments)
	go start()

	time.Sleep(250 * time.Millisecond)

	response, err := http.Get("http://localhost:8008/v1/")
	assert.NoError(err)
	if err == nil {
		defer response.Body.Close()
		assert.Equal(200, response.StatusCode)
		assert.Equal("application/json", response.Header.Get("Content-Type"))
		body, err := ioutil.ReadAll(response.Body)
		assert.NoError(err)
		assert.Equal("{\"message\":\"Welcome to gobstopper you are connected to database: golang\"}", string(body))
	}

}
