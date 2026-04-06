package main

import (
	"net/http"
	"testing"

	"greenlight.isez.dev/internal/assert"
)

func TestHealthcheckHandler(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/v1/healthcheck", "")

	assert.Equal(t, code, http.StatusOK)
	assert.StringContains(t, body, "available")
	assert.StringContains(t, body, "development")
}

func TestListMoviesHandler(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		auth     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid Request",
			urlPath:  "/v1/movies",
			wantCode: http.StatusOK,
			auth:     "",
			wantBody: "Moana",
		},
		{
			name:     "Invalid sort parameter",
			urlPath:  "/v1/movies?sort=invalid",
			wantCode: http.StatusUnprocessableEntity,
			auth:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath, tt.auth)
			assert.Equal(t, code, tt.wantCode)
			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}
