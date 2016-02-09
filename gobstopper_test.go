package gobstopper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServerWithoutDB(t *testing.T) {
	assert := assert.New(t)

	conf := Config{
		Port:       8000,
		PathPrefix: "/v1/prefix",
	}

	server, err := NewServer(conf)
	assert.NoError(err)
	assert.NotNil(server)
	assert.IsType(&Server{}, server)
	assert.Equal(8000, server.Port)
	assert.Equal("/v1/prefix", server.PathPrefix)
	assert.NotNil(server.Router)
	assert.Empty(server.Middleware)
	assert.Nil(server.DB)
}

func TestNewServerErrorsOnConnectionTypeFailure(t *testing.T) {
	assert := assert.New(t)

	conf := Config{
		Port:       8000,
		PathPrefix: "/v1/prefix",
		Connection: "bad connection",
	}

	server, err := NewServer(conf)
	assert.Error(err)
	assert.Nil(server)
}

func TestServerErrorsOnConnectionFailure(t *testing.T) {
	assert := assert.New(t)

	conf := Config{
		Port:       8000,
		PathPrefix: "/v1/prefix",
		Connection: "mysql|bad connection",
	}

	server, err := NewServer(conf)
	assert.Error(err)
	assert.Nil(server)
}
