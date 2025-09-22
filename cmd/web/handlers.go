package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"letsGo/internal/models"
	"letsGo/internal/validator"
	"log"
	"net/http"
	"strings"
)

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

var albums []models.Album

func (app *application) usersList(w http.ResponseWriter, r *http.Request) {
	users, err := app.users.List()
	if err != nil {
		log.Print(err.Error())
		app.serverError(w, r, err)
		return
	}
	data := app.newTemplateData(r)
	data.Users = users

	app.render(w, r, http.StatusOK, "userslist.html", data)
}
func (app *application) employeesList(w http.ResponseWriter, r *http.Request) {
	employees, err := app.employees.List()
	if err != nil {
		log.Print(err.Error())
		app.serverError(w, r, err)
		return
	}
	data := app.newTemplateData(r)
	data.Employees = employees

	app.render(w, r, http.StatusOK, "employees.html", data)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Albums = albums
	app.render(w, r, http.StatusOK, "search.html", data)
}

func (app *application) search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if len(query) < 2 {
		json.NewEncoder(w).Encode([]string{})
		return
	}

	var matches []string
	lowerQuery := strings.ToLower(query)
	for _, city := range dict {
		if strings.HasPrefix(strings.ToLower(city), lowerQuery) {
			matches = append(matches, city)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches)
}

func (app *application) submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	artist := r.PostFormValue("artist")
	artistId, _ := findKeyByValue(dataMap, artist)
	log.Printf("Search for artist: %s and ArtistID: %s\n", artist, artistId)
	http.Redirect(w, r, "/", http.StatusSeeOther)

	albums, _ = app.albums.ListAlbums(artistId)
	fmt.Printf("Albums for artist %v: \n", albums)
}

func findKeyByValue(m map[string]string, targetValue string) (string, bool) {
	for key, value := range m {
		if value == targetValue {
			return key, true
		}
	}
	return "", false // No match found
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.html", data)
}
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.html", data)
}
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Do some validation checks on the form. We check that both email and
	// password are provided, and also check the format of the email address as
	// a UX-nicety (in case the user makes a typo).
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
		return
	}
	// Check whether the credentials are valid. If they're not, add a generic
	// non-field error message and redisplay the login page.
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	// Use the RenewToken() method on the current session to change the session
	// ID. It's good practice to generate a new session ID when the
	// authentication state or privilege levels change for the user (e.g. login
	// and logout operations).
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Add the ID of the current user to the session, so that they are now
	// 'logged in'.
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	//email := "maria@g.com"
	x := strings.Split(form.Email, "@")
	name := strings.ToUpper(x[0])
	app.sessionManager.Put(r.Context(), "flash", `Welcome `+name)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	// Declare a zero-valued instance of our userSignupForm struct.
	var form userSignupForm
	// Parse the form data into the userSignupForm struct.
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Validate the form contents using our helper functions.
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 4), "password", "This field must be at least 4 characters long")
	// If there are any errors, redisplay the signup form along with a 422
	// status code.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Otherwise add a confirmation flash message to the session confirming that
	// their signup worked.
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	// And redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// Use the RenewToken() method on the current session to change the session
	// ID again.
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Remove the authenticatedUserID from the session data so that the user is
	// 'logged out'.
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	// Add a flash message to the session to confirm to the user that they've been
	// logged out.
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	// Redirect the user to the application home page.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
