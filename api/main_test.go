package api_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/m-kuzmin/simple-rest-api/logging"
)

func TestMain(m *testing.M) {
	logging.GlobalLogger = logging.StdLogger{}

	gin.SetMode(gin.TestMode)

	m.Run()
}
