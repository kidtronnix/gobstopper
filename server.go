package gobstopper

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func NewServer(conf Config) (*Server, error) {
	Router := mux.NewRouter().StrictSlash(true)

	if conf.Negroni == nil {
		conf.Negroni = negroni.Classic()
	}

	var DB *sqlx.DB
	if conf.Connection != "" {
		db, err := ConnectToSQL(conf.Connection)
		if err != nil {
			return nil, err
		}
		if err := db.Ping(); err != nil {
			return nil, err
		}
		DB = db
	}

	return &Server{
		&conf,
		Router,
		DB,
	}, nil
}

type ServerInterface interface {
	Route(Route)
	RouteGroup(string, *negroni.Negroni, func(ServerInterface))
}

// Server is the main struct responsible for holding all the required pieces of the stack.
// It consists of a sqlx.DB connection, an array of negroni middleware functions, a Negroni instance, an array of routes handlers, a mux.Router instance,
// the service port to run from and the path prefix.
type Server struct {
	*Config
	Router *mux.Router
	DB     *sqlx.DB
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

func (s *Server) RouteGroup(path string, routing func(*Group)) {
	g := &Group{
		PathPrefix: s.PathPrefix + path,
		Router:     mux.NewRouter(),
		Negroni:    negroni.New(),
	}

	routing(g)
	g.Negroni.Use(negroni.Wrap(g.Router))
	s.Router.PathPrefix(s.PathPrefix + path).Handler(g.Negroni)
}

func (s *Server) Start() {
	s.Negroni.Use(negroni.Wrap(s.Router))
	s.Negroni.Run(fmt.Sprintf(":%d", s.Port))
}

func (s *Server) Middleware(ms ...negroni.Handler) {
	for _, m := range ms {
		s.Negroni.Use(m)
	}
}

func (s *Server) Route(r Route) {
	if s.PathPrefix == "" {
		s.Router.Path(r.Path).Methods(r.Method).Handler(r.Handler)
	} else {
		s.Router.PathPrefix(s.PathPrefix).Path(r.Path).Methods(r.Method).Handler(r.Handler)
	}
}
