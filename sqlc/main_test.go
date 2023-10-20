package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// temporarily defined as constants but will be retrieved from env vars soon
const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:password@localhost:5499/simple_bank?sslmode=disable"
)

// this vars contains all methods for the sqlc queries we defined
var testQueries *Queries

// serves as entry point for tests, test files need to communicate with the DB
func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
