package main

import "net/http"

func (app *application) logError(r *http.Request, err error) {
	method := r.Method
	uri := r.URL.RequestURI()

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (app *application) clientError(w http.ResponseWriter, status int, message string) {
	http.Error(w, message, status)
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	http.Error(w, message, http.StatusInternalServerError)
}
