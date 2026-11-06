package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /passwords/check", app.checkPasswordHandler)
	mux.HandleFunc("POST /passwords", app.addPasswordHandler)

	mux.Handle("GET /", http.FileServer(http.Dir("./static")))

	return mux
}
