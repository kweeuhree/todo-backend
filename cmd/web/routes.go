package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter" // router
	"github.com/justinas/alice"           // middleware
)

func (app *application) routes() http.Handler {
	// Initialize the router.
	router := httprouter.New()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/api/", app.home)
	router.HandlerFunc(http.MethodGet, "/api/todo/view", app.todoView)      // fixed path
	router.HandlerFunc(http.MethodPost, "/api/todo/create", app.todoCreate) // fixed path

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(router)
}
