package api

import (
	"github.com/gin-gonic/gin"
	"github.com/m-kuzmin/simple-rest-api/api/swaggerui"
	"github.com/m-kuzmin/simple-rest-api/logging"
)

func NewGinRouter(server *Server) *gin.Engine {
	router := gin.New()

	swaggerui.Mount(router, "swaggerui") // TODO Move this to main
	router.PUT("/users", server.CreateOrUpdateUsers)

	logging.Infof("Gin router is set-up.")

	return router
}
