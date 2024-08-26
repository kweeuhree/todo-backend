package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid" // router
	"todo-backend.kweeuhree/internal/models"
	"todo-backend.kweeuhree/internal/validator"
)

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

// user authentication routes
// sign up a new user
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "Attempting to create a new user...")
	// declare a zero-valued instance of userInput struct
	var form userSignUpInput

	// parse the form data into the struct
	err := decodeJSON(w, r, &form)
	if err != nil {
		return
	}

	log.Printf("Received new user details: %s", form)

	// Validate the form contents using our helper functions.
	// form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	// form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	// form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	// form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	// form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	form.Validate()
	if !form.Valid() {
		err := encodeJSON(w, http.StatusOK, form.FieldErrors)
		if err != nil {
			// app.serverError(w, err)
			json.NewEncoder(w).Encode(err)
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
	err = encodeJSON(w, http.StatusOK, response)
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
	err := decodeJSON(w, r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	log.Printf("Attempting to authenticate user: %s", form)

	// validate input
	form.Validate()
	if !form.Valid() {
		err := encodeJSON(w, http.StatusOK, form.FieldErrors)
		if err != nil {
			json.NewEncoder(w).Encode(form.FieldErrors)
			app.serverError(w, err)
		}
		return
	}

	// Check whether the credentials are valid. If they're not, add a generic
	// non-field error message and re-display the login page.
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			app.errorLog.Printf("Failed adding user to database: %s", err)
			json.NewEncoder(w).Encode("Invalid credentials")
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

	// Set the flash message
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	app.setFlash(r.Context(), "Login successful!")

	// Create a response that includes both ID and body
	response := UserResponse{
		Uuid:  id,
		Email: form.Email,
		Flash: app.getFlash(r.Context()),
	}

	// Write the response struct to the response as JSON
	err = encodeJSON(w, http.StatusOK, response)
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
	err := decodeJSON(w, r, &form)
	if err != nil {
		json.NewEncoder(w).Encode(err)
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
	err = encodeJSON(w, http.StatusOK, response)
	if err != nil {
		app.serverError(w, err)
		json.NewEncoder(w).Encode(err)
		return
	}

	fmt.Println(w, "Logged out the user")

}
