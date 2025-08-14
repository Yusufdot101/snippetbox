package main

import (
	"fmt"
	"net/http"
	"strconv"
	// "text/template"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	// check if the url path is exactly /
	// this avoids paths like /foo/bar mapping to /
	// if its not return page not found and exist the function
	if r.URL.Path != "/" {
		app.clientError(w, 404)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, apiSuccess{Result: snippets})
	// intialize a slice containing the paths to the two files.
	// the base file path must come first
	// files := []string{
	// 	"./ui/html/base.tmpl.html",
	// 	"./ui/html/partials/nav.tmpl.html",
	// 	"./ui/html/pages/home.tmpl.html",
	// }

	// use template.ParseFiles() function to read the files and store the the
	// templates inot into a template set
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	// use the ExecuteTemplate() method to write the content of the base
	// template sa the response body
	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }
}

func (app *application) ViewSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.clientError(w, 400)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		app.clientError(w, 404)
		return
	}
	WriteJSON(w, http.StatusOK, apiSuccess{Result: snippet})
}

func (app *application) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, 405)
		return
	}

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
