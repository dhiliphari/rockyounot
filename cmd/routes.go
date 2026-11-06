package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/passwords/check", app.checkPasswordHandler)
	mux.HandleFunc("POST /api/passwords", app.addPasswordHandler)

	return mux
}
