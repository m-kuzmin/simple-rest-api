package swaggerui

import (
	"path"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

func Mount(engine *gin.Engine, swaggerBasePath string) {
	if engine == nil {
		engine = gin.Default()
	}

	engine.GET(path.Join(swaggerBasePath, "*any"), ginswagger.WrapHandler(swaggerfiles.Handler))
}
