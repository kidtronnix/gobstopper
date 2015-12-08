package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/seedboxtech/adx/middleware/filters"
	"github.com/seedboxtech/adx/service"
)

var (
	port       int
	prefix     string
	connection string
)

func initFlags(f *flag.FlagSet, arguments []string) {
	f.IntVar(&port, "port", 8000, "Port to listen to")
	f.StringVar(&prefix, "prefix", "/v1/stats", "Prefix of service that appears in URL")
	f.StringVar(&connection, "connection", "mysql|root:root@/golang?charset=utf8", "Connection string to DB")
	f.Parse(arguments)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// To get db connection:
	// db := database.GetConnection(r)
	w.Header().Set("Content-Type", "application/json")
	val := "world"
	fmt.Fprintf(w, "{\"hello\":\"%s\"}", val)
}

func main() {

	initFlags(flag.CommandLine, os.Args[1:])

	service, err := service.NewService(port, prefix, connection)
	if err != nil {
		log.Fatalf("Unable to start service: %s", err)
	}

	service.AddRoute("GET", "/", handler)
	service.AddMiddleware(filters.FilterMiddleware)

	service.Start()
}
