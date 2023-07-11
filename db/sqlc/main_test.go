package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open("postgres", "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable")
	if err != nil {
		log.Fatal("Not able to connect to the database")
	}
	testQueries = New(conn)
	os.Exit(m.Run())
}
