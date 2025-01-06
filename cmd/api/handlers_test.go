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

// func TestListMoviesHandler(t *testing.T) {
// 	app := newTestApplication(t)

// 	ts := newTestServer(t, app.routes())
// 	defer ts.Close()

// 	tests := []struct {
// 		name     string
// 		urlPath  string
// 		wantCode int
// 		wantBody string
// 	}{
// 		{
// 			name:     "Valid Request",
// 			urlPath:  "/v1/movies",
// 			wantCode: http.StatusOK,
// 			wantBody: `[
//         {
//             "id": 2,
//             "title": "Moana",
//             "year": 2016,
//             "runtime": "107 mins",
//             "genres": [
//                 "animation",
//                 "adventure"
//             ],
//             "version": 1
//         }
//     	]`,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			code, _, body := ts.get(t, tt.urlPath)

// 			assert.Equal(t, code, tt.wantCode)

// 			if tt.wantBody != "" {
// 				assert.StringContains(t, body, tt.wantBody)
// 			}
// 		})
// 	}
// }
