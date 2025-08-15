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

func (app *application) viewSnippet(w http.ResponseWriter, r *http.Request) {
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

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
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

func (app *application) createSnippetPost(w http.ResponseWriter, r *http.Request) {
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
	http.Redirect(w, r, fmt.Sprintf("/snippets/view/%d", id), http.StatusSeeOther)
}

type apiError struct {
	Error any `json:"error"`
}

type apiSuccess struct {
	Result any `json:"result"`
}

type apiFunction func(http.ResponseWriter, *http.Request) error
