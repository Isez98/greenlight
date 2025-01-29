package main

import (
	"bytes"
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"testing"

	"greenlight.isez.dev/internal/data"
)

func newTestApplication(t *testing.T, db_dsn string) *application {
	db_conn, err := newTestDB(t, db_dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db_conn.Close()

	return &application{
		config: config{
			env: "development",
		},
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		models: data.NewModels(db_conn),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string, authToken string) (int, http.Header, string) {
	// rs, err := ts.Client().Get(ts.URL + urlPath)
	req, err := http.NewRequest("GET", ts.URL+urlPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+authToken)

	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

func newTestDB(t *testing.T, db_dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", db_dsn)
	if err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("../../internal/data/testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	t.Cleanup(func() {
		defer db.Close()

		script, err := os.ReadFile("../../internal/data/testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
	})

	return db, err
}
