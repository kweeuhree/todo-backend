package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"todo-backend.kweeuhree/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/". If it
	// doesn't, use the http.NotFound() function to send a 404 response to the client.
	// return from the handler. Failing to return the handler would result
	// in "hello world" message being printed as well
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("hello world"))
}

func (app *application) todoView(w http.ResponseWriter, r *http.Request) {
	// get id from URL query
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	// return a 404 Not Found in case of invalid id or error
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	todo, err := app.todos.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// write the todo data as a plain-text HTTP response body
	fmt.Fprintf(w, "%+v", todo)
}
