package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// models
	"todo-backend.kweeuhree/internal/models"

	// environment variables
	"github.com/joho/godotenv"

	// we need the driver’s init() function to run so that it can register itself with the
	// database/sql package. The trick to getting around this is to alias the package name
	// to the blank identifier. This is standard practice for most of Go’s SQL drivers
	_ "github.com/go-sql-driver/mysql" // with underscore

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
)

// Define an application struct to hold the application-wide dependencies for
// the web application.
type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	users          *models.UserModel
	todos          *models.TodoModel
	sessionManager *scs.SessionManager
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

	// Create a new MySQL session store using the connection pool.
	store := mysqlstore.New(db)
	// Initialize a new session manager.
	sessionManager := scs.New()
	// Use the MySQL session store with the session manager.
	sessionManager.Store = store
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		users:          &models.UserModel{DB: db},
		todos:          &models.TodoModel{DB: db},
		sessionManager: sessionManager,
	}

	// Initialize a tls.Config struct to hold the non-default TLS settings we
	// want the server to use. In this case the only thing that we're
	// changing is the curve preferences value, so that only elliptic curves with
	// assembly implementations are used.
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:      *addr,
		ErrorLog:  errorLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
		// connection timeouts
		// -- all keep-alive connections will be automatically closed
		// -- after 1 minute of inactivity
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		// -- prevent the data that the handler returns
		// -- from taking too long to write
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)

	// ListenAndServeTLS() starts HTTPS server
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
