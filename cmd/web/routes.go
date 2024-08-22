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

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/api", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/api/todo/view", dynamic.ThenFunc(app.todoView))                      // fixed path
	router.Handler(http.MethodPost, "/api/todo/create", dynamic.ThenFunc(app.todoCreate))                 // fixed path
	router.Handler(http.MethodPut, "/api/todo/update/:id", dynamic.ThenFunc(app.todoUpdate))              // fixed path
	router.Handler(http.MethodPut, "/api/todo/toggle-status/:id", dynamic.ThenFunc(app.todoToggleStatus)) // fixed path
	router.Handler(http.MethodDelete, "/api/todo/delete/:id", dynamic.ThenFunc(app.todoDelete))           // fixed path
	// router.Handler(http.MethodGet, "/api/test-cookie", dynamic.ThenFunc(app.testCookie))

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(router)
}
