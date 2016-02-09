package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/seedboxtech/gobstopper"
	"github.com/seedboxtech/gobstopper/examples/middleware"
)

func main() {
	conf := gobstopper.Config{
		Port: 8000,
		// PathPrefix: "/v1",
	}

	server, err := gobstopper.NewServer(conf)
	if err != nil {
		log.Fatal("Unable to start server!")
	}

	server.Route(gobstopper.Route{
		Method: "GET",
		Path:   "/foo",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "hello")
		}),
	})

	n := negroni.New(negroni.HandlerFunc(middleware.TerribleAuthMiddleware))
	server.RouteGroup("/admin", n, func(r *mux.Router) {
		r.Path("/admin").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "admin")
		})
	})

	server.Start()

}
