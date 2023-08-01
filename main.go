//nolint:wsl // main() looks better this way
package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	httpReadTimeout = time.Minute
)

func main() {
	logging.GlobalLogger = logging.StdLogger{}

	logging.Infof("Connecting to Postgres")
	postgres := MustSetupPostgres()
	logging.Infof("Connected to Postgres")

	server := api.NewServer(postgres)

	gin.SetMode(gin.ReleaseMode)
	router := api.NewGinRouter(server)

	httpServer := StartServer(router)
	logging.Infof("Server started")

	WaitForCtrcC()
	logging.Infof("Shutting down the server")

	err := httpServer.Shutdown(context.Background())
	if err != nil {
		logging.Fatalf("Error during shutdown: %s", err)
	}
	logging.Infof("[1/2] HTTP handler stopped")

	if err = postgres.Close(); err != nil {
		logging.Errorf("Error closing PostgreSQL connection: %s", err)
	}
	logging.Infof("[2/2] SQL connection closed")
	logging.Infof("Server gracefully shut down")
}

func StartServer(engine http.Handler) *http.Server {
	server := &http.Server{
		Addr:        bindToPort,
		Handler:     engine,
		ReadTimeout: httpReadTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logging.Errorf("HTTP server error: %s", err)
		}

		logging.Infof("HTTP server shutdown")
	}()

	return server
}

func WaitForCtrcC() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt
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
