package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"greenlight.isez.dev/internal/data"
	"greenlight.isez.dev/internal/validator"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredenttailsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredenttailsResponse(w, r)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) verifyTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Vary", "Authorization")

	authorizationHeader := r.Header.Get("Authorization")

	if authorizationHeader == "" {
		r = app.contextSetUser(r, data.AnonymousUser)
		app.writeJSON(w, http.StatusOK, envelope{"valid": false}, nil)
		return
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		app.invalidCredenttailsResponse(w, r)
		return
	}

	token := headerParts[1]

	v := validator.New()

	data.ValidateTokenPlaintext(v, token)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	valid, err := app.models.Tokens.ValidateToken(token)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"valid": valid}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
