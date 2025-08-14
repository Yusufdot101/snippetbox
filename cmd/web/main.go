package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Yusufdot101/snippetbox/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
}

func main() {
	// the dafult port if the addr flag is not set
	defaultPort := ":4000"
	defaultDSN := "web:REMOVED_PASSWORD@/snippetbox?parseTime=true"
	// a cammmand line flag named "addr", with default value defaultPort
	// short text explaining what the flag controls
	// store in addr variable
	addr := flag.String("addr", defaultPort, "HTTP newtwork address")
	dsn := flag.String("dsn", defaultDSN, "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
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
