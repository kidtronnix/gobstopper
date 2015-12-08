package service

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/seedboxtech/adx/db"
)

// Service is the main struct responsible for holding all the required pieces of the stack.
// It consists of a sqlx.DB connection, an array of negroni middleware functions, an array of routes handlers and a mux.Router.
// Before start up of the service the routes are added to the mux.Router, the DB connection is added to
// the context of the request so that connection can be used in any function that has access to the *http.Request object.
// The request will also pass through all the middleware added to the service in the order it was added.
type Service struct {
	port       int
	prefix     string
	Router     *mux.Router
	Negroni    *negroni.Negroni
	routes     []Route
	middleware []Middleware
	DB         *sqlx.DB
}

// Route is the struct for routes. It consists of the HTTP method ie GET, POST, etc...
// The path is the URL path to the route, this path will be added to the service prefix.
// The handler is the function responsible for actually handling the request.
type Route struct {
	path    string
	method  string
	handler RouteHandler
}

// RouteHandler is simply an alias of the http.HandlerFunc.
// It is the function used to handle requests and write responses.
type RouteHandler http.HandlerFunc

// Middleware is a type that has the required signature to act as negroni Middleware.
// A new middleware should call the next function to yeild to the next middleware in the chain.
type Middleware func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

// NewService makes a new instance of a service. The port defines the port the service will run from.
// The prefix is added as a path before all routes. The connection is a slightly modified sql connection string
// in the form of `driverName|dataSourceName` example, mysql|root:root@/golang?charset=utf8.
// If you do not need SQL DB connection simply leave the string empty.
// To access the *sqlx.DB connection in route handlers or middleware functions call database.GetFromContext(r *http.Request)
func NewService(port int, prefix string, connection string) (*Service, error) {
	s := &Service{}
	s.port = port
	s.prefix = prefix
	s.Router = mux.NewRouter()
	s.routes = []Route{}
	s.middleware = []Middleware{}
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
func (s *Service) DBConnectionMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	database.SetConnectionInRequestContext(r, s.DB)
	next(rw, r)
}

// Start is a convenience function used to start the service once all the routes and middleware has been added.
// On calling it loops through all the routes adds them to the mux.Router
// makes a classic negroni, adds middleware stack and routes to that negroni instance and
// starts the service.
func (s *Service) Start() {
	s.Negroni = negroni.Classic()

	for _, route := range s.routes {
		s.Router.HandleFunc(s.prefix+route.path, route.handler).Methods(route.method)
	}

	// s.Negroni.Use(negroni.HandlerFunc(s.DBConnectionMiddleware))

	for _, middleware := range s.middleware {
		s.Negroni.Use(negroni.HandlerFunc(middleware))
	}

	s.Negroni.UseHandler(s.Router)
	s.Negroni.Run(fmt.Sprintf(":%d", s.port))
}

// AddRoute will add a route to the service.
func (s *Service) AddRoute(method, path string, handler RouteHandler) {
	r := Route{
		method:  method,
		path:    path,
		handler: handler,
	}
	s.routes = append(s.routes, r)
}

// AddMiddleware will add middleware to the http stack.
func (s *Service) AddMiddleware(m Middleware) {
	s.middleware = append(s.middleware, m)
}
