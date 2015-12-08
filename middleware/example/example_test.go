package example

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://example.com?token=123", nil)
	assert.NoError(err)
	next := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	Middleware(w, req, next)
	assert.Equal(200, w.Code)
}

func TestMiddlewareFail(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://example.com?token=letmeinguys", nil)
	assert.NoError(err)
	next := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	Middleware(w, req, next)
	assert.Equal(401, w.Code)
	assert.Equal("Authentication failed", w.Body.String())
}
