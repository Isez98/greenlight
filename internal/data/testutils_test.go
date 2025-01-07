package data

import (
	"database/sql"
	"flag"
	"os"
	"testing"
)

var db_dsn string

func TestMain(m *testing.M) {
	flag.StringVar(&db_dsn, "db-dsn", os.Getenv("TEST_DB_DSN"), "Test DSN")
	flag.Parse()
	m.Run()
}

func newTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", db_dsn)
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
