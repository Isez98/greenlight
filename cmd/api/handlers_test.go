package main

import (
	"flag"
	"net/http"
	"os"
	"testing"

	"greenlight.isez.dev/internal/assert"
)

var db_dsn string

func TestMain(m *testing.M) {
	flag.StringVar(&db_dsn, "db-dsn", os.Getenv("TEST_DB_DSN"), "Test DSN")
	flag.Parse()
	m.Run()
}

// func TestHealthcheckHandler(t *testing.T) {
// 	app := newTestApplication(t, db_dsn)

// 	ts := newTestServer(t, app.routes())
// 	defer ts.Close()

// 	healtcheckTest := struct {
// 		status      string
// 		environment string
// 		version     string
// 	}{
// 		status:      "available",
// 		environment: "development",
// 		version:     "1.0.0",
// 	}

// 	code, _, body := ts.get(t, "/v1/healthcheck", "")

// 	assert.Equal(t, code, http.StatusOK)
// 	assert.StringContains(t, body, healtcheckTest.status)
// 	assert.StringContains(t, body, healtcheckTest.version)
// 	assert.StringContains(t, body, healtcheckTest.environment)
// }

func TestListMoviesHandler(t *testing.T) {
	app := newTestApplication(t, db_dsn)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	authToken := ""

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
			auth:     authToken,
			wantBody: `[
        {
            "id": 1,
            "title": "Black Panther",
            "year": 2018,
            "runtime": "134 mins",
            "genres": [
                "action",
                "adventure"
            ],
            "version": 1
        }
    	]`,
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
