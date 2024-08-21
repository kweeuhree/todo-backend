package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"      // validation
	"unicode/utf8" // validation

	// validation
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter" // router
	"todo-backend.kweeuhree/internal/models"
)

// Input struct for creating and updating todos
type TodoInput struct {
	Body string `json:"body"`
}

// Response struct for returning todo data
type TodoResponse struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	todos, err := app.todos.All()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Set the Content-Type header to application/json if you are sending JSON
	w.Header().Set("Content-Type", "application/json")

	// Write the todos to the response as JSON
	err = json.NewEncoder(w).Encode(todos)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

// read
func (app *application) todoView(w http.ResponseWriter, r *http.Request) {
	// use the ParamsFromContext() function to retrieve a slice containing

	// these parameter names and values like so:
	params := httprouter.ParamsFromContext(r.Context())

	// We can then use the ByName() method to get the value of the "id" named
	// parameter from the slice and validate it as normal.
	id, err := strconv.Atoi(params.ByName("id"))

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

// create
func (app *application) todoCreate(w http.ResponseWriter, r *http.Request) {
	// check if method is POST
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")

		// use clientError helper instead of a http.Error shortcut
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// get body from the request, and define a new id with UUID
	// err := r.ParseForm()

	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }

	// Decode the JSON body into the input struct
	var input TodoInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Initialize a map to hold any validation errors for the form fields.
	fieldErrors := make(map[string]string)

	// Check that the Content value isn't blank.
	if strings.TrimSpace(input.Body) == "" {
		fieldErrors["body"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(input.Body) > 200 {
		fieldErrors["body"] = "This field cannot be more than 100 characters long"
	}

	// If there are any errors, dump them in a plain text HTTP response and
	// return from the handler.
	if len(fieldErrors) > 0 {
		fmt.Fprint(w, fieldErrors)
		return
	}

	// Set the Content-Type header to application/json if you are sending JSON
	w.Header().Set("Content-Type", "application/json")

	newId := uuid.New().String()

	// Insert the new todo using the ID and body
	id, err := app.todos.Insert(newId, input.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Create a response that includes both ID and body
	response := TodoResponse{
		ID:   id,
		Body: input.Body,
	}

	// Write the response struct to the response as JSON
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	// http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}

// update
func (app *application) todoUpdate(w http.ResponseWriter, r *http.Request) {

}

// delete
func (app *application) todoDelete(w http.ResponseWriter, r *http.Request) {

}
