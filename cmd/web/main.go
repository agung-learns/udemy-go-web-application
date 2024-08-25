package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"html/template"
	"learn1/internal/driver"
	"learn1/internal/models"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"
const cssVersion = "1"

var session *scs.SessionManager

type config struct {
	port int
	env  string
	api  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
}

type application struct {
	config        config
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	version       string
	DB            models.DBModel
	Session       *scs.SessionManager
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}
	app.infoLog.Printf("Starting server on port %d (%s)", app.config.port, app.config.env)
	return srv.ListenAndServe()
}

func main() {
	gob.Register(TransactionData{})
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "HTTP server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|production)")
	flag.StringVar(&cfg.api, "api", "http://localhost:4002", "API URL")
	flag.StringVar(&cfg.db.dsn, "dsn", "postgres://golearn_user:golearn_password@localhost:5432/golearn_db?sslmode=disable", "Database DSN")
	flag.Parse()

	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	conn, err := driver.OpenDB(cfg.db.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer conn.Close()

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Store = postgresstore.New(conn)

	tc := make(map[string]*template.Template)

	app := &application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		version:       version,
		DB: models.DBModel{
			DB: conn,
		},
		Session: session,
	}

	if err := app.serve(); err != nil {
		app.errorLog.Fatal(err)
	}
}
