package data

import (
	"database/sql"
	"flag"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	var db_dsn string
	flag.StringVar(&db_dsn, "db-dsn", os.Getenv("TEST_DB_DSN"), "Test DSN")
	flag.Parse()
	os.Exit(m.Run())
}

func newTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("TEST_DB_DSN"))
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
