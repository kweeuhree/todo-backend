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

	// uprotected application routes using the "dynamic" middleware chain, use nosurf middleware
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)

	// todo routes
	router.Handler(http.MethodGet, "/api", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/api/todo/view", dynamic.ThenFunc(app.todoView))
	// csrf token route
	router.Handler(http.MethodGet, "/api/csrf-token", dynamic.ThenFunc(app.CSRFToken))
	// test
	// router.Handler(http.MethodGet, "/api/test-cookie", dynamic.ThenFunc(app.testCookie))

	// user routes
	// create a new user
	router.Handler(http.MethodPost, "/api/user/signup", dynamic.ThenFunc(app.userSignup))
	// authenticate and login the user
	router.Handler(http.MethodPost, "/api/user/login", dynamic.ThenFunc(app.userLogin))

	// protected application routes, which uses requireAuthentication middleware
	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodPost, "/api/todo/create", protected.ThenFunc(app.todoCreate)) // fixed path
	router.Handler(http.MethodPut, "/api/todo/update/:id", protected.ThenFunc(app.todoUpdate))
	router.Handler(http.MethodPut, "/api/todo/toggle-status/:id", protected.ThenFunc(app.todoToggleStatus))
	router.Handler(http.MethodDelete, "/api/todo/delete/:id", protected.ThenFunc(app.todoDelete))
	// logout the user
	router.Handler(http.MethodPost, "/api/user/logout", protected.ThenFunc(app.userLogout))
	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(router)
}
