package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	// models
	"todo-backend.kweeuhree/internal/models"

	// environment variables
	"github.com/joho/godotenv"

	// we need the driver’s init() function to run so that it can register itself with the
	// database/sql package. The trick to getting around this is to alias the package name
	// to the blank identifier. This is standard practice for most of Go’s SQL drivers
	_ "github.com/go-sql-driver/mysql" // with underscore
)

// Define an application struct to hold the application-wide dependencies for
// the web application.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	todos    *models.TodoModel
}

func main() {
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// DSN string with loaded env variables
	DSNstring := fmt.Sprintf("%s:%s@/%s?parseTime=true", dbUser, dbPassword, dbName)

	// define  new command-line flag for the mysql dsn string
	dsn := flag.String("dsn", DSNstring, "MySQL data source name")
	addr := flag.String("addr", ":4000", "HTTP network address")

	// parse flags
	flag.Parse()

	// error and info logs
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// create connection pool, pass openDB() the dsn from the command-line flag
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// defer a call to db.Close() so that the connection pool is closed before
	// the main() function exits
	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		todos:    &models.TodoModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)

	// use assignment operator as the err variable is already declared above
	err = srv.ListenAndServe()
	// in case of errors log and exit
	errorLog.Fatal(err)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool for a given dsn
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
