package swaggerui

import (
	"path"

	"github.com/gin-gonic/gin"
	"github.com/m-kuzmin/simple-rest-api/logging"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

func Mount(engine *gin.Engine, swaggerBasePath string) {
	if engine == nil {
		engine = gin.Default()
	}

	engine.GET(path.Join(swaggerBasePath, "*any"), ginswagger.WrapHandler(swaggerfiles.Handler))
	logging.Infof("Swagger available at: %s", path.Join(swaggerBasePath, "index.html"))
}
