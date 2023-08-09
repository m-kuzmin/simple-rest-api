package api

import (
	"github.com/gin-gonic/gin"
	"github.com/m-kuzmin/simple-rest-api/logging"
)

func NewGinRouter(server *Server) *gin.Engine {
	router := gin.New()

	router.POST("/users", server.CreateUsersBody)
	router.POST("/users/upload", server.CreateUsersUpload)
	router.GET("/users/search", server.SearchUsers)

	logging.Infof("Gin router is set-up.")

	return router
}
