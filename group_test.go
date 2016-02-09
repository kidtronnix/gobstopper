package gobstopper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGroupMiddleware(t *testing.T) {
	assert := assert.New(t)

	g := Group{
		PathPrefix: "",
		Negroni:    negroni.New(),
		Router:     mux.NewRouter(),
	}

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "http://example.com?token=123", nil)
	assert.NoError(err)

	MiddlewareTestKey := "TESTKEY"
	next := func(w http.ResponseWriter, r *http.Request) {
		a := context.Get(r, MiddlewareTestKey)
		assert.Equal("Set by middleware", a)
	}

	m := negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		context.Set(r, MiddlewareTestKey, "Set by middleware")
		next(w, r)
	})

	g.Middleware(m)

	handlers := g.Negroni.Handlers()
	assert.Len(handlers, 1)

	if len(handlers) == 1 {
		handlers[0].ServeHTTP(w, r, next)
	}

}

func TestGroupRoute(t *testing.T) {
	assert := assert.New(t)

	g := Group{
		PathPrefix: "/v1",
		Negroni:    negroni.New(),
		Router:     mux.NewRouter(),
	}

	g.Route(Route{
		Method: "GET",
		Path:   "/endpoint",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "foo")
		}),
	})

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "http://example.com:8000/v1/endpoint", nil)
	assert.NoError(err)

	g.Router.ServeHTTP(w, r)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("foo", w.Body.String())

}
