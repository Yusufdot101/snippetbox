package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/Yusufdot101/snippetbox/internal/models"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// load the env vars
	err := godotenv.Load()
	if err != nil {
		errorLog.Fatal("Error loading .env file")
	}
	// the default values for the flags if not set on run time
	defaultPort := ":4000"
	dbPassowrd := os.Getenv("DB_PASSWORD")
	defaultDSN := "web:" + dbPassowrd + "@/snippetbox?parseTime=true"
	// a cammmand line flag named addr or dsn, with default value
	// short text explaining what the flag controls
	// store in appropriate variable
	addr := flag.String("addr", defaultPort, "HTTP newtwork address")
	dsn := flag.String("dsn", defaultDSN, "MySQL data source name")
	flag.Parse()

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	// initialize a new http.Server struct. we set the Addr and Handler fields so
	// that the server uses the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom errorLog logger in
	// the event of any problems.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// the value returned from the flag.String() functiona is a pointer to the flag
	// value, not the value itselt
	// write the messages using custom loggers
	infoLog.Printf("Server listening on port: %s", *addr)

	// Call the ListenAndServe() method on our new http.Server struct.
	err = srv.ListenAndServe()

	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
