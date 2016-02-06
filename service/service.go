// Package service helps make services with a mux, negroni, sqlx stack.
package service

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/seedboxtech/gobstopper/db"
)

// Service is the main struct responsible for holding all the required pieces of the stack.
// It consists of a sqlx.DB connection, an array of negroni middleware functions, a Negroni instance, an array of routes handlers, a mux.Router instance,
// the service port to run from and the path prefix.
type Service struct {
	Port        int
	Prefix      string
	router      *mux.Router
	Router      *mux.Router
	RouteGroups []*RouteGroup
	Middleware  []Middleware
	DB          *sqlx.DB
}

// Route is the struct for routes. It consists of the HTTP method ie GET, POST, etc...
// The path is the URL path to the route, this path will be added to the service prefix.
// The handler is the function responsible for actually handling the request.
type Route struct {
	Path    string
	Method  string
	Handler http.Handler
}

type RouteGroup struct {
	Path       string
	Router     *mux.Router
	Middleware []Middleware
}

// Middleware is a type that has the required signature to act as negroni Middleware.
// A new middleware should call the next function to yeild to the next middleware in the chain.
type Middleware func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

// NewService makes a new instance of a service. The port defines the port the service will run from.
// The prefix is added as a path before all routes. The connection is a slightly modified sql connection string
// in the form of `driverName|dataSourceName` example, mysql|root:root@/golang?charset=utf8.
// If you do not need SQL DB connection simply leave the string empty.
func NewService(port int, prefix string, connection string) (*Service, error) {
	s := &Service{}
	s.Port = port
	s.Prefix = prefix
	s.router = mux.NewRouter().StrictSlash(true)
	s.Router = s.router.PathPrefix(prefix).Subrouter().StrictSlash(true)
	s.Middleware = []Middleware{}
	if connection != "" {
		db, err := database.ConnectToSQL(connection)
		if err != nil {
			return nil, err
		}
		if err := db.Ping(); err != nil {
			return nil, err
		}
		s.DB = db
	}

	return s, nil
}

// DBConnectionMiddleware is a middleware function that sets the *sqlx.DB connection into
// the request context.
// To access the *sqlx.DB connection in route handlers or middleware functions call database.GetFromContext(r *http.Request)
func (s *Service) DBConnectionMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	database.SetConnectionInRequestContext(r, s.DB)
	next(rw, r)
}

func (s *Service) AddMiddlewareToRouteGroup(rg *RouteGroup) {
	n := negroni.New()
	addMiddlewareToRouter(n, rg.Middleware, rg.Router)
	s.Router.PathPrefix(rg.Path).Handler(n)
}

func addMiddlewareToRouter(n *negroni.Negroni, ms []Middleware, r *mux.Router) {
	for _, m := range ms {
		n.Use(negroni.HandlerFunc(m))
	}
	n.Use(negroni.Wrap(r))
}

func (s *Service) NewRouteGroup(path string) *RouteGroup {
	router := mux.NewRouter()
	rg := RouteGroup{
		Router: router.PathPrefix(s.Prefix + path).Subrouter(),
		Path:   path,
	}
	s.RouteGroups = append(s.RouteGroups, &rg)
	return &rg
}

// Start is a convenience function used to start the service once all the routes and middleware have been added.
// On calling it loops through all the routes adds them to the mux.Router
// makes a classic negroni, adds middleware stack and routes to that negroni instance and
// starts the service.
func (s *Service) Start() {
	n := negroni.Classic()

	for _, rg := range s.RouteGroups {
		s.AddMiddlewareToRouteGroup(rg)
	}

	addMiddlewareToRouter(n, s.Middleware, s.router)

	n.Run(fmt.Sprintf(":%d", s.Port))
}

// AddRouteHandler will add a route to service.Routes for a given http.Handler.
func (s *Service) AddRouteHandler(method, path string, handler http.Handler) {
	addRouteHandler(s.Router, method, path, handler)
}

// AddRouteHandlerFunc will add a route to service.Routes for a given http.HandlerFunc.
func (s *Service) AddRouteHandlerFunc(method, path string, handler http.HandlerFunc) {
	addRouteHandlerFunc(s.Router, method, path, handler)
}

// AddMiddleware will add middleware to service.Middleware.
func (s *Service) AddMiddleware(m Middleware) {
	s.Middleware = append(s.Middleware, m)
}

func (r *RouteGroup) AddRouteHandler(method, path string, handler http.Handler) {
	addRouteHandler(r.Router, method, path, handler)
}

func (r *RouteGroup) AddRouteHandlerFunc(method, path string, handler http.HandlerFunc) {
	addRouteHandlerFunc(r.Router, method, path, handler)
}

func (r *RouteGroup) AddMiddleware(m Middleware) {
	r.Middleware = append(r.Middleware, m)
}

func addRouteHandlerFunc(r *mux.Router, method, path string, handler http.HandlerFunc) {
	r.HandleFunc(path, handler).Methods(method).Name(method + " " + path)
}

func addRouteHandler(r *mux.Router, method, path string, handler http.Handler) {
	r.Handle(path, handler).Methods(method).Name(method + " " + path)
}
