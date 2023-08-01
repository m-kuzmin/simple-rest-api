package sqlc_test

import (
	"database/sql"
	"testing"

	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/m-kuzmin/simple-rest-api/db"
	"github.com/m-kuzmin/simple-rest-api/db/sqlc"
	"github.com/m-kuzmin/simple-rest-api/logging"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/users?sslmode=disable"
)

var testQueries *sqlc.Queries //nolint:gochecknoglobals // This is used by all tests to connect to the DB.

func TestMain(m *testing.M) {
	logging.GlobalLogger = logging.StdLogger{}

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		logging.Fatalf("Cannot connect to database: %s", err)
	}

	if err = db.PostgresMigrateUp(conn, "file://../migrations", "users"); err != nil {
		logging.Fatalf("Failed to migrate test db: %s", err)
	}

	testQueries = sqlc.New(conn)

	m.Run()
}
