package logger

import (
	"context"
)

// ==================== 全局实例 ====================

var (
	std Logger
)

func Init(opts ...Option) {
	std = New(opts...).WithCaller(2)
}

// Default 获取全局 Logger
func Default() Logger {
	return std
}

// ==================== 快捷方法 ====================

func Debug(ctx context.Context, msg string, fields ...Field) { Ctx(ctx).Log(DEBUG, msg, fields...) }
func Info(ctx context.Context, msg string, fields ...Field)  { Ctx(ctx).Log(INFO, msg, fields...) }
func Warn(ctx context.Context, msg string, fields ...Field)  { Ctx(ctx).Log(WARN, msg, fields...) }
func Error(ctx context.Context, msg string, fields ...Field) { Ctx(ctx).Log(ERROR, msg, fields...) }
func Fatal(ctx context.Context, msg string, fields ...Field) { Ctx(ctx).Log(FATAL, msg, fields...) }
func Errw(ctx context.Context, err error, fields ...Field) {
	Ctx(ctx).Log(ERROR, err.Error(), fields...)
}
