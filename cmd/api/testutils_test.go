package main

import (
	"bytes"
	"database/sql"
	"flag"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"greenlight.isez.dev/internal/data/mocks"
)

func newTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("TEST_GREENLIGHT_DB"))
	if err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("./testdata/setup.sql")
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

		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
	})

	return db
}

func newTestApplication(t *testing.T) *application {
	var cfg config
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	db, err := openDB(cfg)
	if err != nil {
		os.Exit(1)
	}
	defer db.Close()

	return &application{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		config: cfg,
		models: mocks.TestModels_Mock(db),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)
	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
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
