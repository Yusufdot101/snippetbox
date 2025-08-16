package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Yusufdot101/snippetbox/internal/models"
	"github.com/Yusufdot101/snippetbox/internal/validator"
	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// use template.ParseFiles() function to read the files and store the the
	// templates inot into a template set
	data := app.newTemplateData(r)
	data.Snippets = snippets

	page := "home.tmpl.html"
	app.render(w, http.StatusOK, page, data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// retrieve a slice containing the paramaters in the url
	params := httprouter.ParamsFromContext(r.Context())

	// get the value of "id" parameter
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.clientError(w, 400)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, 404)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	page := "view.tmpl.html"
	app.render(w, http.StatusOK, page, data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	page := "create.tmpl.html"
	app.render(w, http.StatusOK, page, data)
}

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	if !app.isAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var form snippetCreateForm

	err = app.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	permittedExpiresValues := []string{"1", "7", "365"}
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PremittedInt(form.Expires, permittedExpiresValues...), "expires", "This field must be in ["+strings.Join(permittedExpiresValues, ", ")+"]")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		page := "create.tmpl.html"
		app.render(w, http.StatusBadRequest, page, data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully!")

	http.Redirect(w, r, fmt.Sprintf("/snippets/view/%d", id), http.StatusSeeOther)
}

type userCreateForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userCreateForm{}
	page := "signup.tmpl.html"
	app.render(w, http.StatusOK, page, data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	var form userCreateForm

	err = app.formDecoder.Decode(&form, r.Form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(
		validator.NotBlank(form.Name),
		"name",
		"This field cannot be blank",
	)
	form.CheckField(
		validator.NotBlank(form.Email),
		"email",
		"This field cannot be blank",
	)
	form.CheckField(
		validator.NotBlank(form.Password),
		"password",
		"This field cannot be blank",
	)

	form.CheckField(
		validator.MinChars(form.Password, 8),
		"password",
		"This must be at least 8 characters long",
	)

	form.CheckField(
		validator.Matches(form.Email, validator.EmailRX),
		"email",
		"This field must be a vaild email address",
	)

	if !form.Valid() {
		page := "signup.tmpl.html"
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusBadRequest, page, data)
		return
	}
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			page := "signup.tmpl.html"
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusBadRequest, page, data)
			return
		}
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	page := "login.tmpl.html"
	app.render(w, http.StatusOK, page, data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	var form userLoginForm
	err = app.formDecoder.Decode(&form, r.Form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(
		validator.NotBlank(form.Email),
		"email", "This field cannot be blank",
	)
	form.CheckField(
		validator.Matches(form.Email, validator.EmailRX),
		"email", "This field must be a valid email address",
	)
	form.CheckField(
		validator.NotBlank(form.Password),
		"password", "This field cannot be blank",
	)

	if !form.Valid() {
		page := "login.tmpl.html"
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusBadRequest, page, data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvaildCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			page := "login.tmpl.html"
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusBadRequest, page, data)
			return
		}
		app.serverError(w, err)
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/snippets/create", http.StatusSeeOther)

}
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type apiError struct {
	Error any `json:"error"`
}

type apiSuccess struct {
	Result any `json:"result"`
}
