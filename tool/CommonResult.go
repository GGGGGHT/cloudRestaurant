package tool

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	FAILED = iota - 1
	SUCCESS
)

func Success(ctx *gin.Context, v interface{}) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": SUCCESS,
		"msg":  "success",
		"data": v,
	})
}

func Failed(ctx *gin.Context, v interface{}) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": FAILED,
		"msg":  "error",
		"data": v,
	})
}