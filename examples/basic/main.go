// A highly opinionated golang web stack; mux, negroni, sqlx. It consists of two core packages, 'service' and 'db'. See 'main.go' for an example of how to use these packages.
// There is also an example middleware included to show how to add custom middleware to your stack.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	_ "github.com/go-sql-driver/mysql"
	"github.com/seedboxtech/gobstopper"
	"github.com/seedboxtech/gobstopper/examples/middleware"
)

var conf gobstopper.Config

// Example request handler
func handler(w http.ResponseWriter, r *http.Request) {
	// Get db connection for request context
	db := gobstopper.GetConnection(r)
	name, _ := gobstopper.GetDatabaseName(db)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"message\":\"Welcome to gobstopper you are connected to database: %s\"}", name)
}

func initFlags(f *flag.FlagSet, arguments []string) {
	f.IntVar(&conf.Port, "port", 8000, "Port to listen to")
	f.StringVar(&conf.PathPrefix, "prefix", "", "Prefix of service that appears in URL")
	f.StringVar(&conf.Connection, "connection", "mysql|root:root@/golang?charset=utf8", "Connection string to DB")
	f.Parse(arguments)
}

func start() {
	server, err := gobstopper.NewServer(conf)
	if err != nil {
		log.Fatalf("Unable to start service: %s", err)
	}

	server.Negroni.Use(negroni.HandlerFunc(server.DBConnectionMiddleware))
	server.Negroni.Use(negroni.HandlerFunc(middleware.TerribleAuthMiddleware))

	// Index Route
	server.Route(gobstopper.Route{
		Method:  "GET",
		Path:    "/",
		Handler: http.HandlerFunc(handler),
	})

	// Admin Routes
	adminRoutes := server.NewRouteGroup("/admin")
	adminRoutes.AddMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		fmt.Println("Route specific middleware: only for admins")
		next(w, r)
	})
	adminRoutes.AddRouteHandlerFunc("GET", "/foo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "woop 2 foo")
	})

	// Start server
	server.Start()
}

func main() {
	initFlags(flag.CommandLine, os.Args[1:])
	start()
}