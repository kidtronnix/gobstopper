package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var conn string

func init() {
	DB := os.Getenv("DB")
	switch DB {
	case "travis":
		conn = "mysql|root:@/golang?charset=utf8"
	default:
		conn = "mysql|root:root@/golang?charset=utf8"
	}
}

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

	arguments := []string{"-port", "8008", "-prefix", "/v1", "-connection", conn}
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	initFlags(f, arguments)
	go start()

	time.Sleep(250 * time.Millisecond)

	response, err := http.Get("http://localhost:8008/v1/?token=123")
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
