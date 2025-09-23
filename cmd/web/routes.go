package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.Handle("GET /{$}", app.sessionManager.LoadAndSave(http.HandlerFunc(app.home)))

	mux.Handle("GET /search", app.sessionManager.LoadAndSave(http.HandlerFunc(app.search)))

	mux.Handle("POST /submit", app.sessionManager.LoadAndSave(http.HandlerFunc(app.submitHandler)))

	mux.Handle("GET /records/{id}", app.sessionManager.LoadAndSave(http.HandlerFunc(app.recordHandler)))

	mux.Handle("GET /employees/list", app.sessionManager.LoadAndSave(http.HandlerFunc(app.employeesList)))

	mux.Handle("GET /user/signup", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userSignup)))

	mux.Handle("GET /user/login", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userLogin)))

	mux.Handle("POST /user/login", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userLoginPost)))

	mux.Handle("POST /user/signup", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userSignupPost)))

	mux.Handle("POST /user/logout", app.sessionManager.LoadAndSave(app.requireAuthentication(http.HandlerFunc(app.userLogoutPost))))

	mux.Handle("GET /user/list", app.sessionManager.LoadAndSave(app.requireAuthentication(http.HandlerFunc(app.usersList))))

	return app.recoverPanic(app.logRequest(commonHeaders(mux)))
}
