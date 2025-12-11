package base

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

func OK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Resp{
		Code: 0,
		Msg:  "ok",
	})
}

func OKWithData(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Resp{
		Code: 0,
		Msg:  "ok",
		Data: data,
	})
}

func Fail(ctx *gin.Context, code int, msg string) {
	ctx.AbortWithStatusJSON(http.StatusOK, Resp{
		Code: code,
		Msg:  msg,
	})
}

func FailWithErr(ctx *gin.Context, code int, err error) {
	ctx.AbortWithStatusJSON(http.StatusOK, Resp{
		Code: code,
		Msg:  err.Error(),
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
	ctx.AbortWithStatusJSON(statusCode, Resp{
		Code: statusCode,
		Msg:  err.Error(),
	})
}
