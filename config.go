package gobstopper

import (
	"net/http"

	"github.com/codegangsta/negroni"
)

type Config struct {
	Port            int                                      `REQUIRED`
	PathPrefix      string                                   `OPTIONAL`
	Connection      string                                   `OPTIONAL`
	Negroni         *negroni.Negroni                         `OPTIONAL DEFAULT TO CLASSIC`
	NotFoundHandler func(http.ResponseWriter, *http.Request) `OPTIONAL DEFAULT TO BUILTIN MUX`
}

type Route struct {
	Method  string
	Path    string
	Handler http.Handler
}
