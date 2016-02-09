package gobstopper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/jmoiron/sqlx"
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
	assert.NotNil(server.Negroni)
	assert.Nil(server.DB)
}

func TestNewServerWitDB(t *testing.T) {
	assert := assert.New(t)

	conf := Config{
		Port:       8000,
		PathPrefix: "/v1/prefix",
		Connection: "mysql|" + conn,
	}

	server, err := NewServer(conf)
	assert.NoError(err)
	assert.NotNil(server)
	assert.NotNil(server.DB)
	err = server.DB.Ping()
	assert.NoError(err)

}

func TestNewServerErrorsOnConnectionTypeFailure(t *testing.T) {
	assert := assert.New(t)

	conf := Config{
		Port:       8000,
		PathPrefix: "/v1/prefix",
		Connection: "bad connection type|root:root",
	}

	service, err := NewServer(conf)
	assert.Error(err)
	assert.Nil(service)
}

func TestNewServerErrorsOnConnectionFailure(t *testing.T) {
	assert := assert.New(t)

	conf := Config{
		Port:       8000,
		PathPrefix: "/v1/prefix",
		Connection: "mysql|bad connection",
	}

	service, err := NewServer(conf)
	assert.Error(err)
	assert.Nil(service)
}

func TestServerDBConnectionMiddleware(t *testing.T) {
	assert := assert.New(t)

	conf := Config{
		Port:       8000,
		PathPrefix: "/v1/prefix",
		Connection: "mysql|" + conn,
	}

	server, err := NewServer(conf)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "http://example.com?token=123", nil)
	assert.NoError(err)
	next := func(w http.ResponseWriter, r *http.Request) {
		db := context.Get(r, DatabaseKey).(*sqlx.DB)
		assert.Equal(server.DB, db)
	}

	server.DBConnectionMiddleware(w, r, next)
}

func TestServerMiddleware(t *testing.T) {
	assert := assert.New(t)

	conf := Config{
		Port:       8000,
		PathPrefix: "/v1/prefix",
		Negroni:    negroni.New(),
	}

	server, err := NewServer(conf)

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

	server.Middleware(m)

	handlers := server.Negroni.Handlers()
	assert.Len(handlers, 1)

	if len(handlers) == 1 {
		handlers[0].ServeHTTP(w, r, next)
	}

}

func TestServerRoute(t *testing.T) {
	assert := assert.New(t)

	conf := Config{
		Port:       8000,
		PathPrefix: "/v1/prefix",
		Connection: "mysql|" + conn,
	}

	server, err := NewServer(conf)

	route := Route{
		Method: "GET",
		Path:   "/endpoint",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "foo")
		}),
	}

	server.Route(route)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "http://example.com:8000/v1/prefix/endpoint", nil)
	assert.NoError(err)

	server.Router.ServeHTTP(w, r)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("foo", w.Body.String())

}

func TestServerRouteGroup(t *testing.T) {
	assert := assert.New(t)

	conf := Config{
		Port:       8000,
		PathPrefix: "/v1/prefix",
		Connection: "mysql|" + conn,
	}

	AdminMiddleware := negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		token := r.FormValue("token")
		if token == "123" {
			next(rw, r)
			return
		}
		rw.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(rw, "Authentication failed")
	})

	server, err := NewServer(conf)

	server.Route(Route{
		Method: "GET",
		Path:   "/foo",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "foo")
		}),
	})

	server.RouteGroup("/admin", func(group *Group) {
		group.Middleware(AdminMiddleware)
		group.Route(Route{
			Method: "GET",
			Path:   "",
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, "admin")
			}),
		})
	})

	// Admin route should be blocked without token
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "http://example.com:8000/v1/prefix/admin", nil)
	assert.NoError(err)
	server.Router.ServeHTTP(w, r)
	assert.Equal(http.StatusUnauthorized, w.Code)

	// Admin route should be OK with correct token
	w = httptest.NewRecorder()
	r, err = http.NewRequest("GET", "http://example.com:8000/v1/prefix/admin?token=123", nil)
	assert.NoError(err)
	server.Router.ServeHTTP(w, r)
	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("admin", w.Body.String())

	// Index route should not have AdminMiddleware attached
	w = httptest.NewRecorder()
	r, err = http.NewRequest("GET", "http://example.com:8000/v1/prefix/foo", nil)
	assert.NoError(err)
	server.Router.ServeHTTP(w, r)
	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("foo", w.Body.String())
}

func TestServerStarts(t *testing.T) {
	assert := assert.New(t)

	conf := Config{
		Port:       8000,
		PathPrefix: "/v1",
		Connection: "mysql|" + conn,
		Negroni:    negroni.New(),
	}
	server, err := NewServer(conf)
	assert.NoError(err)

	server.Middleware(negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		fmt.Fprint(rw, "middleware... ")
		next(rw, r)
	}))

	server.Route(Route{
		Method: "GET",
		Path:   "",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "handler...")
		}),
	})

	// Start service and wait for it start
	go server.Start()
	time.Sleep(250 * time.Millisecond)

	response, err := http.Get("http://localhost:8000/v1")
	assert.NoError(err)
	if err == nil {
		defer response.Body.Close()
		assert.Equal(200, response.StatusCode)
		body, err := ioutil.ReadAll(response.Body)
		assert.NoError(err)
		assert.Equal("middleware... handler...", string(body))
	}

}
