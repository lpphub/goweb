package logx

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lpphub/goweb/pkg/logger"
)

const (
	// TraceLogIDHeader 请求中携带的 Trace Log ID Header
	TraceLogIDHeader = "X-Trace-LogId"
)

func init() {
	// 兼容处理gin.context
	logger.RegisterCtxAdapter(func(ctx context.Context) context.Context {
		if gCtx, ok := ctx.(*gin.Context); ok {
			return gCtx.Request.Context()
		}
		return ctx
	})
}

type accessLogConfig struct {
	skipPaths map[string]struct{}
}

type AccessLogOption func(*accessLogConfig)

func defaultConfig() *accessLogConfig {
	return &accessLogConfig{
		skipPaths: make(map[string]struct{}),
	}
}

// SkipPaths 跳过指定的完整路径，如 /health
func SkipPaths(paths ...string) AccessLogOption {
	return func(cfg *accessLogConfig) {
		for _, p := range paths {
			if p != "" {
				cfg.skipPaths[p] = struct{}{}
			}
		}
	}
}

// GinAccessLog Gin 请求访问日志中间件（支持跳过路径）
func GinAccessLog(opts ...AccessLogOption) gin.HandlerFunc {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if _, ok := cfg.skipPaths[path]; ok {
			c.Next()
			return
		}

		// 注入 trace log id 到 context
		ctx := logger.CtxWithField(c.Request.Context(), logger.Str("logId", resolveTraceLogID(c)))
		c.Request = c.Request.WithContext(ctx)

		logger.Info(ctx, "gin access",
			logger.Str("path", fmt.Sprintf("[%s %s]",
				c.Request.Method, c.Request.RequestURI,
			)),
		)

		c.Next()
	}
}

func resolveTraceLogID(c *gin.Context) string {
	if c.Request != nil {
		if logID := c.GetHeader(TraceLogIDHeader); logID != "" {
			return logID
		}
	}
	return GenerateTraceLogID()
}

func GenerateTraceLogID() string {
	return strconv.FormatUint(uint64(time.Now().UnixNano())&0x7FFFFFFF|0x80000000, 10)
}
