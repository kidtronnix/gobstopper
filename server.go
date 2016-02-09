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

	if conf.NotFoundHandler != nil {
		Router.NotFoundHandler = conf.NotFoundHandler
	}

	if conf.Negroni == nil {
		conf.Negroni = negroni.Classic()
	}

	var DB *sqlx.DB
	server := Server{
		&conf,
		Router,
		DB,
	}

	if conf.Connection != "" {
		db, err := ConnectToSQL(conf.Connection)
		if err != nil {
			return nil, err
		}
		if err := db.Ping(); err != nil {
			return nil, err
		}
		server.DB = db
		server.Negroni.Use(negroni.HandlerFunc(server.DBConnectionMiddleware))
	}

	return &server, nil

}

type Server struct {
	*Config
	Router *mux.Router
	DB     *sqlx.DB
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
	s.Router.Path(s.PathPrefix + r.Path).Methods(r.Method).Handler(r.Handler)

}
