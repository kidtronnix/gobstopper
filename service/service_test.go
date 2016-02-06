package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
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

type TestHandler struct{}

func (th TestHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprint(rw, "ok")
}

func TestServiceNewRouteGroup(t *testing.T) {
	assert := assert.New(t)

	s := Service{}
	rg := s.NewRouteGroup("/foo")

	assert.Len(s.RouteGroups, 1)
	assert.Equal(rg.Path, "/foo")
	assert.IsType(&mux.Router{}, rg.Router)
}

func TestServiceAddMiddlewareToRouter(t *testing.T) {
	assert := assert.New(t)

	n := negroni.New()
	m := []Middleware{func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		next(w, r)
	}}
	r := mux.NewRouter()
	h := TestHandler{}
	r.Handle("/foo", h)

	addMiddlewareToRouter(n, m, r)

	nh := n.Handlers()

	assert.Len(nh, 2)

}

func TestRouteGroupAddMiddleware(t *testing.T) {
	assert := assert.New(t)

	rg := RouteGroup{}
	m := func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		next(w, r)
	}
	rg.AddMiddleware(m)

	assert.Len(rg.Middleware, 1)
}

func TestRouteGroupAddRouteHandler(t *testing.T) {
	assert := assert.New(t)

	rg := RouteGroup{
		Router: mux.NewRouter(),
		Path:   "/foo",
	}

	h := TestHandler{}

	rg.AddRouteHandler("GET", "/", h)

	r := rg.Router.Get("GET /")

	assert.Equal(h, r.GetHandler())
}

func TestRouteGroupAddRouteHandlerFunc(t *testing.T) {
	assert := assert.New(t)

	rg := RouteGroup{
		Router: mux.NewRouter(),
		Path:   "/foo",
	}

	h := TestHandler{}

	rg.AddRouteHandlerFunc("GET", "/", h.ServeHTTP)

	r := rg.Router.Get("GET /")

	assert.NotNil(r.GetHandler())
}

func TestServiceAddRouteHandler(t *testing.T) {
	assert := assert.New(t)

	service, err := NewService(8001, "/v1/prefix", conn)
	assert.NoError(err)

	h := TestHandler{}

	service.AddRouteHandler("GET", "/", h)

	r := service.Router.Get("GET /")

	assert.Equal(h, r.GetHandler())
}

func TestServiceAddRouteHandlerFunc(t *testing.T) {
	assert := assert.New(t)

	service, err := NewService(8001, "/v1/prefix", conn)
	assert.NoError(err)

	h := TestHandler{}

	service.AddRouteHandlerFunc("GET", "/", h.ServeHTTP)

	r := service.Router.Get("GET /")

	assert.NotNil(r)
}

func TestServiceStartsWithDB(t *testing.T) {
	assert := assert.New(t)

	service, err := NewService(8000, "/v1/prefix", conn)
	assert.NoError(err)
	assert.NotNil(service)
	assert.IsType(&Service{}, service)
	assert.Equal(8000, service.Port)
	assert.Equal("/v1/prefix", service.Prefix)
	assert.NotNil(service.Router)
	assert.Empty(service.Middleware)
	err = service.DB.Ping()
	assert.NoError(err)

	service.AddRouteHandlerFunc("GET", "/", func(rw http.ResponseWriter, r *http.Request) {
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

func TestServiceStartsWithRouteGroup(t *testing.T) {
	assert := assert.New(t)

	service, err := NewService(8001, "/v1", "")
	assert.NoError(err)

	adminRoutes := service.NewRouteGroup("/admin")
	adminRoutes.AddMiddleware(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		rw.Header().Set("Foo", "bar")
		next(rw, r)
	})
	adminRoutes.AddRouteHandlerFunc("GET", "/foo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "foo")
	})

	service.AddRouteHandlerFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "index")
	})

	go service.Start()

	time.Sleep(250 * time.Millisecond)

	response, err := http.Get("http://localhost:8001/v1/admin/foo")
	assert.NoError(err)
	if err == nil {
		defer response.Body.Close()
		assert.Equal("bar", response.Header.Get("Foo"))
		assert.Equal(200, response.StatusCode)
		body, err := ioutil.ReadAll(response.Body)
		assert.NoError(err)
		assert.Equal("foo\n", string(body))
	}

	response, err = http.Get("http://localhost:8001/v1/")
	assert.NoError(err)
	if err == nil {
		defer response.Body.Close()
		assert.Equal("", response.Header.Get("Foo"))
	}
}

func TestNewServiceWithoutDB(t *testing.T) {
	assert := assert.New(t)

	service, err := NewService(8000, "/v1/prefix", "")
	assert.NoError(err)
	assert.NotNil(service)
	assert.IsType(&Service{}, service)
	assert.Equal(8000, service.Port)
	assert.Equal("/v1/prefix", service.Prefix)
	assert.NotNil(service.Router)
	assert.Empty(service.Middleware)
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
