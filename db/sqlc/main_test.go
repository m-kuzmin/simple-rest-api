package sqlc_test

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/lib/pq"
	"github.com/m-kuzmin/simple-rest-api/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/users?sslmode=disable"
)

var testQueries *sqlc.Queries //nolint:gochecknoglobals // This is used by all tests to connect to the DB.

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to database:", err) //nolint:forbidigo // Fine in tests
	}

	testQueries = sqlc.New(conn)

	m.Run()
}
