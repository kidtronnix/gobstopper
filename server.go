package gobstopper

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// type Server interface {
// 	Route(Route)
// 	RouteGroup(string, *negroni.Negroni, func(Server))
// }

// Server is the main struct responsible for holding all the required pieces of the stack.
// It consists of a sqlx.DB connection, an array of negroni middleware functions, a Negroni instance, an array of routes handlers, a mux.Router instance,
// the service port to run from and the path prefix.
type Server struct {
	*Config
	Router  *mux.Router
	Negroni *negroni.Negroni
	DB      *sqlx.DB
}

type Route struct {
	Method  string
	Path    string
	Handler http.Handler
}

func (s *Server) DBConnectionMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	SetConnectionInRequestContext(r, s.DB)
	next(rw, r)
}

func (s *Server) RouteGroup(path string, n *negroni.Negroni, routing func(*mux.Router)) {
	r := mux.NewRouter()
	routing(r)
	n.Use(negroni.Wrap(r))
	s.Router.PathPrefix(s.PathPrefix + path).Handler(n)
}

func (s *Server) Start() {
	s.Negroni.Use(negroni.Wrap(s.Router))
	s.Negroni.Run(fmt.Sprintf(":%d", s.Port))
}

func (s *Server) Route(r Route) {
	if s.PathPrefix == "" {
		s.Router.Path(r.Path).Methods(r.Method).Handler(r.Handler)
	} else {
		s.Router.PathPrefix(s.PathPrefix).Path(r.Path).Methods(r.Method).Handler(r.Handler)
	}
}
