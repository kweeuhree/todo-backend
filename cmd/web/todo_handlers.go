package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	ID    string `json:"id"`
	Body  string `json:"body"`
	Flash string
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
	// Get the value of the "id" named parameter
	params := httprouter.ParamsFromContext(r.Context())
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

	app.sessionManager.Put(r.Context(), "flash", "Todo successfully added!")

	// write the todo data as a plain-text HTTP response body
	fmt.Fprintf(w, "%+v", todo)
}

// create
func (app *application) todoCreate(w http.ResponseWriter, r *http.Request) {

	// Decode the JSON body into the input struct
	var input TodoInput
	err := decodeJSON(w, r, &input)
	if err != nil {
		return
	}

	log.Printf("Received input.Body: %s", input.Body)

	// validate input
	input.Validate()
	if !input.Valid() {
		encodeJSON(w, http.StatusBadRequest, input.FieldErrors)
		return
	}

	newId := uuid.New().String()

	// Insert the new todo using the ID and body
	id, err := app.todos.Insert(newId, input.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.setFlash(r.Context(), "Todo has been created.")

	// Create a response that includes both ID and body
	response := TodoResponse{
		ID:    id,
		Body:  input.Body,
		Flash: app.getFlash(r.Context()),
	}

	// Write the response struct to the response as JSON
	err = encodeJSON(w, http.StatusOK, response)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

// update
func (app *application) todoUpdate(w http.ResponseWriter, r *http.Request) {
	log.Printf("Attempting update...")

	// Get the value of the "id" named parameter
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")
	log.Printf("Current todo id: %s", id)

	if id == "" {
		app.notFound(w)
		log.Printf("Exiting due to invalid id")
		return
	}

	// Decode the JSON body into the input struct
	var input TodoInput
	err := decodeJSON(w, r, &input)
	if err != nil {
		log.Printf("Exiting after decoding attempt...")
		log.Printf("Error message %s", err)
		return
	}

	log.Printf("Received input.Body: %s", input.Body)

	// validate input
	input.Validate()
	if !input.Valid() {
		encodeJSON(w, http.StatusBadRequest, input.FieldErrors)
		return
	}

	// Update the new todo using the ID and body
	err = app.todos.Put(id, input.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.setFlash(r.Context(), "Todo has been updated.")

	// Create a response that includes both ID and body
	response := TodoResponse{
		ID:    id,
		Body:  input.Body,
		Flash: app.getFlash(r.Context()),
	}

	// Write the response struct to the response as JSON
	err = encodeJSON(w, http.StatusOK, response)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) todoToggleStatus(w http.ResponseWriter, r *http.Request) {
	log.Printf("Attempting status toggle...")

	// Get the value of the "id" named parameter
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")
	log.Printf("Current todo id: %s", id)

	if id == "" {
		app.notFound(w)
		log.Printf("Exiting due to invalid id")
		return
	}

	// Delete the todo using the ID
	err := app.todos.Toggle(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

// delete
func (app *application) todoDelete(w http.ResponseWriter, r *http.Request) {
	log.Printf("Attempting deletion...")
	// Set the Content-Type header to application/json
	// w.Header().Set("Content-Type", "application/json")

	// Get the value of the "id" named parameter
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")
	log.Printf("Current todo id: %s", id)

	if id == "" {
		app.notFound(w)
		log.Printf("Exiting due to invalid id")
		return
	}

	// Delete the todo using the ID
	err := app.todos.Delete(id)
	if err != nil {
		app.serverError(w, err)
		return
	} else {
		json.NewEncoder(w).Encode("Deleted successfully!")
		return
	}
}
