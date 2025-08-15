package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Yusufdot101/snippetbox/internal/models"
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
	WriteJSON(w, http.StatusOK, apiSuccess{Result: "Create snippet"})
	return
}

func (app *application) createSnippetPost(w http.ResponseWriter, r *http.Request) {
	title := "first snippet"
	content := "first snippet content"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	WriteJSON(w, http.StatusCreated, apiSuccess{Result: "snippet: " + strconv.Itoa(id) + " created"})
	http.Redirect(w, r, fmt.Sprintf("/snippets/view?id=%d", id), http.StatusSeeOther)
}

type apiError struct {
	Error any `json:"error"`
}

type apiSuccess struct {
	Result any `json:"result"`
}

type apiFunction func(http.ResponseWriter, *http.Request) error
