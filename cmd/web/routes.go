package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.Home)
	mux.HandleFunc("/snippets/view", app.ViewSnippet)
	mux.HandleFunc("/snippets/create", app.CreateSnippet)

	return mux
}
