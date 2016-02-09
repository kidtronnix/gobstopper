package gobstopper

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

func TestServerNewRouteGroup(t *testing.T) {
	assert := assert.New(t)

	s := Server{}
	rg := s.NewRouteGroup("/foo")

	assert.Len(s.RouteGroups, 1)
	assert.Equal(rg.Path, "/foo")
	assert.IsType(&mux.Router{}, rg.Router)
}

func TestServerAddMiddlewareToRouter(t *testing.T) {
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

func TestServerAddRouteHandler(t *testing.T) {
	assert := assert.New(t)

	conf := Config{
		Port:       8001,
		PathPrefix: "/v1/prefix",
		Connection: conn,
	}

	server, err := NewServer(conf)
	assert.NoError(err)

	h := TestHandler{}

	server.AddRouteHandler("GET", "/", h)

	r := server.Router.Get("GET /")

	assert.Equal(h, r.GetHandler())
}

func TestServerAddRouteHandlerFunc(t *testing.T) {
	assert := assert.New(t)

	server, err := NewServer(8001, "/v1/prefix", conn)
	assert.NoError(err)

	h := TestHandler{}

	server.AddRouteHandlerFunc("GET", "/", h.ServeHTTP)

	r := server.Router.Get("GET /")

	assert.NotNil(r)
}

func TestServerStartsWithDB(t *testing.T) {
	assert := assert.New(t)

	server, err := NewServer(8000, "/v1/prefix", conn)
	assert.NoError(err)
	assert.NotNil(server)
	assert.IsType(&Server{}, server)
	assert.Equal(8000, server.Port)
	assert.Equal("/v1/prefix", server.Prefix)
	assert.NotNil(server.Router)
	assert.Empty(server.Middleware)
	err = server.DB.Ping()
	assert.NoError(err)

	server.AddRouteHandlerFunc("GET", "/", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		fmt.Fprint(rw, "ok")
	})

	m := func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		rw.Header().Set("Foo", "bar")
		next(rw, r)
	}
	server.AddMiddleware(server.DBConnectionMiddleware)
	server.AddMiddleware(m)

	go server.Start()

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

func TestServerStartsWithRouteGroup(t *testing.T) {
	assert := assert.New(t)

	server, err := NewServer(8001, "/v1", "")
	assert.NoError(err)

	adminRoutes := server.NewRouteGroup("/admin")
	adminRoutes.AddMiddleware(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		rw.Header().Set("Foo", "bar")
		next(rw, r)
	})
	adminRoutes.AddRouteHandlerFunc("GET", "/foo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "foo")
	})

	server.AddRouteHandlerFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "index")
	})

	go server.Start()

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
