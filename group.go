package gobstopper

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

type Group struct {
	PathPrefix string
	Negroni    *negroni.Negroni
	Router     *mux.Router
}

func (g *Group) Route(r Route) {
	if g.PathPrefix == "" {
		g.Router.Path(r.Path).Methods(r.Method).Handler(r.Handler)
	} else if r.Path == "" {
		g.Router.Path(g.PathPrefix).Methods(r.Method).Handler(r.Handler)
	} else {
		g.Router.PathPrefix(g.PathPrefix).Path(r.Path).Methods(r.Method).Handler(r.Handler)
	}
}

func (g *Group) Middleware(ms ...negroni.Handler) {
	for _, m := range ms {
		g.Negroni.Use(m)
	}
}
