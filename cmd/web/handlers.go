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

// userSignUpInput struct for creating a new user
type userSignUpInput struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form: "-"`
}

type UserResponse struct {
	Uuid  string `json:"uuid"`
	Email string `json:"email"`
	Flash string
}

type userLoginInput struct {
	Email               string `form: "email"`
	Password            string `form: "password"`
	validator.Validator `form: "-" `
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

	// check if method is POST
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")

		// use clientError helper instead of a http.Error shortcut
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Decode the JSON body into the input struct
	var input TodoInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	log.Printf("Received input.Body: %s", input.Body)

	// validate input
	input.CheckField(validator.NotBlank(input.Body), "body", "This field cannot be blank")
	input.CheckField(validator.MaxChars(input.Body, 200), "body", "This field cannot be more than 200 characters long")

	if !input.Valid() {
		err := json.NewEncoder(w).Encode(input.FieldErrors)
		if err != nil {
			app.serverError(w, err)
		}
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

	app.sessionManager.Put(r.Context(), "flash", "This is a flash message!")

	// Write the response struct to the response as JSON
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

// update
func (app *application) todoUpdate(w http.ResponseWriter, r *http.Request) {
	log.Printf("Attempting update...")
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

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
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		log.Printf("Exiting after decoding attempt...")
		log.Printf("Error message %s", err)
		return
	}

	log.Printf("Received input.Body: %s", input.Body)

	// validate input
	input.CheckField(validator.NotBlank(input.Body), "body", "This field cannot be blank")
	input.CheckField(validator.MaxChars(input.Body, 200), "body", "This field cannot be more than 200 characters long")

	if !input.Valid() {
		err := json.NewEncoder(w).Encode(input.FieldErrors)
		if err != nil {
			app.serverError(w, err)
		}
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
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) todoToggleStatus(w http.ResponseWriter, r *http.Request) {
	log.Printf("Attempting status toggle...")
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

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
	w.Header().Set("Content-Type", "application/json")

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
	}
}

// user authentication routes
// sign up a new user
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "Attempting to create a new user...")
	// declare a zero-valued instance of userInput struct
	var form userSignUpInput

	// parse the form data into the struct
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	log.Printf("Received new user details: %s", form)

	// Validate the form contents using our helper functions.
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		err := json.NewEncoder(w).Encode(form.FieldErrors)
		if err != nil {
			app.serverError(w, err)
		}
		return
	}

	newId := uuid.New().String()

	// Try to create a new user record in the database. If the email already
	// exists then add an error message to the form and re-display it.
	err = app.users.Insert(newId, form.Name, form.Email, form.Password)

	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			app.errorLog.Printf("Failed adding user to database: %s", err)
			json.NewEncoder(w).Encode(form.FieldErrors)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.setFlash(r.Context(), "Your signup was successful. Please log in.")

	// Create a response that includes both ID and body
	response := UserResponse{
		Uuid:  newId,
		Email: form.Email,
		Flash: app.getFlash(r.Context()),
	}

	// Write the response struct to the response as JSON
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		app.serverError(w, err)
		return
	}

	fmt.Println(w, "Created a new user...")
}

// authenticate and login the user
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "Attempting to authenticate and login the user...")

	// Decode the form data into the userLoginForm struct.
	var form userLoginInput
	// parse the form data into the struct
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	log.Printf("Attempting to authenticate user: %s", form)

	// checks that email and password are provided
	// and also check the format of the email address as

	// a UX-nicety (in case the user makes a typo).
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		err := json.NewEncoder(w).Encode(form.FieldErrors)
		if err != nil {
			app.serverError(w, err)
		}
		return
	}

	// Check whether the credentials are valid. If they're not, add a generic
	// non-field error message and re-display the login page.
	id, err := app.users.Authenticate(form.Email, form.Password)

	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddNonFieldError("Email or password is incorrect")
			app.errorLog.Printf("Failed adding user to database: %s", err)
			json.NewEncoder(w).Encode(form.FieldErrors)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Use the RenewToken() method on the current session to change the
	// session ID. It's good practice to generate a new session ID when the
	// authentication state or privilege levels changes for the user (e.g.
	// login and logout operations).
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Log the user in
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	// Set the flash message
	app.setFlash(r.Context(), "Login successful!")

	// Create a response that includes both ID and body
	response := UserResponse{
		Uuid:  id,
		Email: form.Email,
		Flash: app.getFlash(r.Context()),
	}

	// Write the response struct to the response as JSON
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		app.serverError(w, err)
		return
	}

	fmt.Println(w, "Authenticated and logged user with ID %d", id)
}

// logout the user
func (app *application) userLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "Attempting to logout the user...")

	// Decode the form data into the userLoginForm struct.
	var form userLoginInput
	// parse the form data into the struct
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// change session ID
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	// remove authenticatedUserID from the session data so that the user is logged out
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.setFlash(r.Context(), "You've been logged out successfully!")

	// Create a response that includes both ID and body
	response := UserResponse{
		Email: form.Email,
		Flash: app.getFlash(r.Context()),
	}

	// Write the response struct to the response as JSON
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		app.serverError(w, err)
		return
	}

	fmt.Println(w, "Logged out the user")

}
