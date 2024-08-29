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
	fmt.Println("Attempting to authenticate and login the user...")

	// Decode the form data into the userLoginInput struct
	var form userLoginInput
	if err := decodeJSON(w, r, &form); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	log.Printf("Attempting to authenticate user: %s", form)

	// Validate input
	form.Validate()
	if !form.Valid() {
		if err := encodeJSON(w, http.StatusOK, form.FieldErrors); err != nil {
			app.serverError(w, err)
			return
		}
		return
	}

	// Check credentials
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			if encodeErr := encodeJSON(w, http.StatusUnauthorized, form.FieldErrors); encodeErr != nil {
				app.serverError(w, encodeErr)
				return
			}
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Renew session token
	if err := app.sessionManager.RenewToken(r.Context()); err != nil {
		app.serverError(w, err)
		return
	}

	// Set flash message
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	app.setFlash(r.Context(), "Login successful!")

	// Create response
	response := UserResponse{
		Uuid:  id,
		Email: form.Email,
		Flash: app.getFlash(r.Context()),
	}

	// Write response
	if err := encodeJSON(w, http.StatusOK, response); err != nil {
		app.serverError(w, err)
		return
	}

	fmt.Printf("Authenticated and logged user with ID %d\n", id)
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
