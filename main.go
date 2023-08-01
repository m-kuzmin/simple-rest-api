//nolint:wsl // main() looks better this way
package main

import (
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/m-kuzmin/simple-rest-api/api"
	"github.com/m-kuzmin/simple-rest-api/db"
	"github.com/m-kuzmin/simple-rest-api/logging"
)

const (
	dbDriver  = "postgres"
	dbAddress = "postgresql://root:secret@postgres:5432/users?sslmode=disable"
	dbName    = "users"

	dbPingRetries      = 5
	dbPingIntervalSecs = 7 // Guestimated amount of time PSQL needs to startup

	migrationsSource = "file:///migrations"

	bindToPort = ":8000"
)

func main() {
	logging.GlobalLogger = logging.StdLogger{}

	logging.Infof("Connecting to Postgres")
	postgres := MustSetupPostgres()
	logging.Infof("Connected to Postgres")

	server := api.NewServer(postgres)

	gin.SetMode(gin.ReleaseMode)
	router := api.NewGinRouter(server)

	logging.Infof("Server started")
	logging.Fatalf("Gin router exit status: %s", router.Run(bindToPort))
}

func MustSetupPostgres() *db.Postgres {
	conn, err := db.ConnectToDBWithRetry(dbDriver, dbAddress, dbPingRetries, dbPingIntervalSecs*time.Second)
	if err != nil {
		logging.Fatalf("failed to connect to PostgreSQL: %s", err)
	}

	if err = db.PostgresMigrateUp(conn, migrationsSource, dbName); err != nil {
		logging.Fatalf("failed to migrate PostgreSQL to latest version: %s", err)
	}

	postgres := db.NewPostgres(conn)

	return postgres
}
