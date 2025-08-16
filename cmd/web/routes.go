package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// Create a new middleware chain containing the middleware specific to our
	// dynamic application routes.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippets/view/:id", dynamic.ThenFunc(app.snippetView))

	router.Handler(http.MethodGet, "/users/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/users/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/users/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/users/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/snippets/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippets/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/users/logout", protected.ThenFunc(app.userLogoutPost))

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeader)

	return standard.Then(router)
}
