package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/m-kuzmin/simple-rest-api/api"
	"github.com/m-kuzmin/simple-rest-api/db"
	"github.com/m-kuzmin/simple-rest-api/logging"
)

const (
	postgresAddr = "postgresql://root:secret@postgres:5432/users?sslmode=disable"
	bindToPort   = ":8000"
)

func main() {
	logging.GlobalLogger = logging.StdLogger{}

	logging.Infof("Connecting to Postgres")

	postgres, err := db.NewPostgres(postgresAddr)
	if err != nil {
		logging.Fatalf("Error connecting to database: %s", err)
	}

	logging.Infof("Connected to Postgres")

	server := api.NewServer(postgres)

	gin.SetMode(gin.ReleaseMode)

	router := api.NewGinRouter(server)

	logging.Infof("Server started")
	logging.Fatalf("Gin router exit status: %s", router.Run(bindToPort))
}
