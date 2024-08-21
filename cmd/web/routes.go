package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter" // router
	"github.com/justinas/alice"           // middleware
)

func (app *application) routes() http.Handler {
	// Initialize the router.
	router := httprouter.New()

	// Create a handler function which wraps our notFound() helper
	// Assign it as the custom handler for 404 Not Found responses
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/api/", app.home)
	router.HandlerFunc(http.MethodGet, "/api/todo/view", app.todoView)        // fixed path
	router.HandlerFunc(http.MethodPost, "/api/todo/create", app.todoCreate)   // fixed path
	router.HandlerFunc(http.MethodPut, "/api/todo/update", app.todoUpdate)    // fixed path
	router.HandlerFunc(http.MethodDelete, "/api/todo/delete", app.todoDelete) // fixed path

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(router)
}
