package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverError helper writes an error messaeg and stack trace to the errorLog,
// then sends a generic 500 internal server errror response to the user
func (app *application) serverError(w http.ResponseWriter, error error) {
	trace := fmt.Sprintf("%s\n%s", error.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	WriteJSON(w, http.StatusInternalServerError, "internal server error")
}

// clientError helper sends a specif status code and corresponding desciption
// to the user
func (app *application) clientError(w http.ResponseWriter, statusCode int) {
	WriteJSON(w, statusCode, http.StatusText(statusCode))
}

func WriteJSON(w http.ResponseWriter, statusCode int, message any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(message)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	//use buffer to see if the rendering works before writing to the ResponseWriter
	buf := new(bytes.Buffer)

	// if there is error return error
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// no errors, its safe to write to ResponseWriter
	w.WriteHeader(status)
	buf.WriteTo(w)
}
