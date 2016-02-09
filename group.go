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
	g.Router.Path(g.PathPrefix + r.Path).Methods(r.Method).Handler(r.Handler)
}

func (g *Group) Middleware(ms ...negroni.Handler) {
	for _, m := range ms {
		g.Negroni.Use(m)
	}
}
