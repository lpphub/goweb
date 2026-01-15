package base

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func OK(ctx *gin.Context) {
	OKWithData(ctx, nil)
}

func OKWithData(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Result{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

func Fail(ctx *gin.Context, code int, msg string) {
	ctx.AbortWithStatusJSON(http.StatusOK, Result{
		Code:    code,
		Message: msg,
	})
}

func FailWithErr(ctx *gin.Context, code int, err error) {
	ctx.AbortWithStatusJSON(http.StatusOK, Result{
		Code:    code,
		Message: err.Error(),
	})
}

func FailWithError(ctx *gin.Context, err error) {
	var e Error
	if ok := errors.As(err, &e); ok {
		Fail(ctx, e.Code, e.Error())
		return
	}
	Fail(ctx, 500, err.Error())
}

func FailWithStatus(ctx *gin.Context, statusCode int, err error) {
	ctx.AbortWithStatusJSON(statusCode, Result{
		Code:    statusCode,
		Message: err.Error(),
	})
}
