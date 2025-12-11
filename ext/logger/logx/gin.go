package logx

import (
	"context"
	"fmt"
	"goweb/ext/logger"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	LogIDHeader = "X-Trace-logId"
)

func init() {
	logger.RegisterCtxExtractor(func(ctx context.Context) context.Context {
		if gCtx, ok := ctx.(*gin.Context); ok {
			return gCtx.Request.Context()
		}
		return ctx
	})
}

func GinLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := logger.CtxWithField(c.Request.Context(), logger.Str("logId", getLogIDFromGin(c)))
		c.Request = c.Request.WithContext(ctx)

		logger.Info(ctx, "gin request",
			logger.Str("path", fmt.Sprintf("[%s %s]", c.Request.Method, c.Request.RequestURI)))

		c.Next()
	}
}

func getLogIDFromGin(ctx *gin.Context) string {
	var logId string
	if ctx.Request != nil && ctx.Request.Header != nil {
		logId = ctx.GetHeader(LogIDHeader)
	}

	if logId == "" {
		logId = GenerateLogID()
	}
	return logId
}

func GenerateLogID() string {
	return strconv.FormatUint(uint64(time.Now().UnixNano())&0x7FFFFFFF|0x80000000, 10)
}
