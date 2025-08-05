package main

import (
	"flag"
	"log"
	"os"
	"rorodata_backend_task/models/store"
)

const SQLITE_DSN = "./rorodata_backend_task"

type config struct {
	port int
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	logger *log.Logger
	models store.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", SQLITE_DSN, "SQLITE3 DSN")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	var (
		m   store.Models
		err error
	)

	m, err = store.New(cfg.db.dsn)
	if err != nil {
		logger.Fatal(err)
	}

	defer m.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: m,
	}

	err = app.serve()
	if err != nil {
		logger.Fatal(err)
	}
}
