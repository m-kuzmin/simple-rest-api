package main

import (
	"github.com/gin-gonic/gin"
	"github.com/m-kuzmin/simple-rest-api/api"
	"github.com/m-kuzmin/simple-rest-api/db"
	"github.com/m-kuzmin/simple-rest-api/logging"
)

func main() {
	logging.GlobalLogger = logging.StdLogger{}

	db := db.NewInMemoryDB()
	server := api.NewServer(db)

	gin.SetMode(gin.ReleaseMode)

	router := api.NewGinRouter(server)

	logging.Fatalf("Gin router exit status: %s", router.Run(":8000"))
}
