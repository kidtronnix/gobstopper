package main

import (
	"flag"
	"testing"

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
