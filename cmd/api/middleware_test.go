package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"greenlight.isez.dev/internal/assert"
)

func TestAuthenticate(t *testing.T) {
	app := newTestApplication(t)

	// A simple next handler that writes the authenticated user's ID to the response.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		fmt.Fprintln(w, user.ID)
	})

	tests := []struct {
		name       string
		authHeader string
		wantCode   int
	}{
		{
			name:       "No auth header sets anonymous user",
			authHeader: "",
			wantCode:   http.StatusOK,
		},
		{
			name:       "Invalid format - missing Bearer prefix",
			authHeader: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			wantCode:   http.StatusUnauthorized,
		},
		{
			name:       "Invalid format - wrong scheme",
			authHeader: "Basic ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			wantCode:   http.StatusUnauthorized,
		},
		{
			name:       "Token too short",
			authHeader: "Bearer tooshort",
			wantCode:   http.StatusUnauthorized,
		},
		{
			// The mock always returns a user for any well-formed token.
			name:       "Valid token authenticates user",
			authHeader: "Bearer ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			wantCode:   http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			app.authenticate(next).ServeHTTP(rr, req)

			assert.Equal(t, rr.Code, tt.wantCode)
		})
	}
}
