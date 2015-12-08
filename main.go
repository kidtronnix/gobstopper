package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/seedboxtech/adx/db"
	"github.com/seedboxtech/adx/middleware/filters"
	"github.com/seedboxtech/adx/service"
)

var (
	port       int
	prefix     string
	connection string
)

// Example request handler
func handler(w http.ResponseWriter, r *http.Request) {
	// Get db connection for request context
	db := database.GetConnection(r)
	name, _ := database.GetDatabaseName(db)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"message\":\"Welcome to gobstopper you are connected to database: %s\"}", name)
}

func initFlags(f *flag.FlagSet, arguments []string) {
	f.IntVar(&port, "port", 8000, "Port to listen to")
	f.StringVar(&prefix, "prefix", "", "Prefix of service that appears in URL")
	f.StringVar(&connection, "connection", "mysql|root:root@/golang?charset=utf8", "Connection string to DB")
	f.Parse(arguments)
}

func start() {
	service, err := service.NewService(port, prefix, connection)
	if err != nil {
		log.Fatalf("Unable to start service: %s", err)
	}
	service.AddRoute("GET", "/", handler)
	// Adds db connection to request context
	service.AddMiddleware(service.DBConnectionMiddleware)
	service.AddMiddleware(filters.FilterMiddleware)
	service.Start()
}

func main() {
	initFlags(flag.CommandLine, os.Args[1:])
	start()
}