package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	// environment variables
	"github.com/joho/godotenv"
)

func secureHeaders(next http.Handler) http.Handler {

	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	reactAddress := os.Getenv("REACT_ADDRESS")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note: This is split across multiple lines for readability. You
		// don't need to do this in your own code.
		// Allow all origins (or specify allowed origins)
		w.Header().Set("Access-Control-Allow-Origin", reactAddress)

		// Allow specific HTTP methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Allow specific headers
		// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method,
			r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

// There are two details about this which are worth explaining:
// Setting the Connection: Close header on the response acts as a
// trigger to make Go’s HTTP server automatically close the current
// connection after a response has been sent. It also informs the
// user that the connection will be closed. Note: If the protocol being
// used is HTTP/2, Go will automatically strip the Connection: Close
// header from the response (so it is not malformed) and send a
// GOAWAY frame.
// The value returned by the builtin recover() function has the type
// any, and its underlying type could be string, error, or something
// else — whatever the parameter passed to panic() was. In our
// case, it’s the string "oops! something went wrong". In the code
// above, we normalize this into an error by using the fmt.Errorf()
// function to create a new error object containing the default
// textual representation of the any value, and then pass this error
// to the app.serverError() helper method.

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a
			// panic or not. If there has...
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")
				// Call the app.serverError helper method to return a 500
				// Internal Server response.
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
