package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	// validation
	// validation
	// validation
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter" // router
	"todo-backend.kweeuhree/internal/models"
	"todo-backend.kweeuhree/internal/validator"
)

// Input struct for creating and updating todos
type TodoInput struct {
	Body string `json:"body"`
	validator.Validator
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

	// Decode the JSON body into the input struct
	var input TodoInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// validate input
	input.CheckField(validator.NotBlank(input.Body), "body", "This field cannot be blank")
	input.CheckField(validator.MaxChars(input.Body, 200), "body", "This field cannot be more than 200 characters long")

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
}

// update
func (app *application) todoUpdate(w http.ResponseWriter, r *http.Request) {

}

// delete
func (app *application) todoDelete(w http.ResponseWriter, r *http.Request) {

}
