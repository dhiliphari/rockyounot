package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (app *application) checkPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(input.Password) == "" {
		app.clientError(w, http.StatusBadRequest, "password must be provided")
		return
	}

	exists, err := app.store.Exists(r.Context(), []byte(input.Password))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := struct {
		Exists bool `json:"exists"`
	}{
		Exists: exists,
	}

	js, err := json.Marshal(data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func (app *application) addPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(input.Password) == "" {
		app.clientError(w, http.StatusBadRequest, "password must be provided")
		return
	}

	err = app.store.Add(r.Context(), []any{input.Password})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte{})
}
