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

	healtcheckTest := struct {
		status      string
		environment string
		version     string
	}{
		status:      "available",
		environment: "development",
		version:     "1.0.0",
	}

	code, _, body := ts.get(t, "/v1/healthcheck")

	assert.Equal(t, code, http.StatusOK)
	assert.StringContains(t, body, healtcheckTest.status)
	assert.StringContains(t, body, healtcheckTest.version)
	assert.StringContains(t, body, healtcheckTest.environment)
}
