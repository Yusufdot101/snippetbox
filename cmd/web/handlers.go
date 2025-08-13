package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

func Home(w http.ResponseWriter, r *http.Request) {
	// check if the url path is exactly /
	// this avoids paths like /foo/bar mapping to /
	// if its not return page not found and exist the function
	if r.URL.Path != "/" {
		WriteJSON(w, http.StatusNotFound, "page not found")
		return
	}

	// intialize a slice containing the paths to the two files.
	// the base file path must come first
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	// use template.ParseFiles() function to read the files and store the the
	// templates inot into a template set
	ts, err := template.ParseFiles(files...)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, apiError{Error: "internal server error"})
		return
	}

	// use the ExecuteTemplate() method to write the content of the base
	// template sa the response body
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		WriteJSON(w, http.StatusInternalServerError, apiError{Error: "internal server error"})
	}
}

func ViewSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		WriteJSON(w, http.StatusBadRequest, apiError{Error: "id that is a number more than 1 is required"})
		return
	}
	WriteJSON(w, http.StatusOK, apiSuccess{Result: "view snippet: " + strconv.Itoa(id)})
}

func CreateSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		WriteJSON(w, http.StatusMethodNotAllowed, apiError{Error: "method not allowed"})
		return
	}
	WriteJSON(w, http.StatusCreated, apiSuccess{Result: "snippet created"})
}

func WriteJSON(w http.ResponseWriter, statusCode int, message any) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(message)
}

type apiError struct {
	Error any `json:"error"`
}

type apiSuccess struct {
	Result any `json:"result"`
}

type apiFunction func(http.ResponseWriter, *http.Request) error
