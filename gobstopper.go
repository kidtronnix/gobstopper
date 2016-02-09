package gobstopper

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func NewServer(conf Config) (*Server, error) {
	Router := mux.NewRouter().StrictSlash(true)
	// Make default behaviour but allow to be passed
	Negroni := negroni.Classic()

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
		Negroni,
		DB,
	}, nil
}
