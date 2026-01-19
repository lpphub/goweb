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
	ctx.JSON(http.StatusOK, Result{
		Code:    0,
		Message: "ok",
	})
}

func OKWithData(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Result{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

func FailWithErr(ctx *gin.Context, err *Error) {
	ctx.AbortWithStatusJSON(http.StatusOK, Result{
		Code:    err.Code,
		Message: err.Error(),
	})
}

func FailWithError(ctx *gin.Context, err error) {
	var bizErr *Error
	if ok := errors.As(err, &bizErr); ok {
		FailWithErr(ctx, bizErr)
		return
	}

	FailWithStatus(ctx, http.StatusInternalServerError, err)
}

func FailWithStatus(ctx *gin.Context, statusCode int, err error) {
	ctx.AbortWithStatusJSON(statusCode, Result{
		Code:    statusCode,
		Message: err.Error(),
	})
}
