package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func errorResponse(ctx *gin.Context, httpCode int, err string) {
	ctx.JSON(httpCode, gin.H{
		"ok":    false,
		"error": err,
	})
}

func errorResponsef(ctx *gin.Context, httpCode int, fmtStr string, a ...any) {
	errorResponse(ctx, httpCode, fmt.Sprintf(fmtStr, a...))
}

func okResponse(ctx *gin.Context, httpCode int) {
	ctx.JSON(httpCode, gin.H{"ok": true})
}
