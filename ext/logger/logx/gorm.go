package logx

import (
	"context"
	"errors"
	"fmt"
	"golib/ext/logger"
	"time"

	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
)

// GormLogger 自定义GORM日志记录器
type GormLogger struct {
	logLevel      glog.LogLevel
	slowThreshold time.Duration
}

// NewGormLogger 创建新的GORM日志记录器
func NewGormLogger() glog.Interface {
	return &GormLogger{
		logLevel: glog.Info,
	}
}

// LogMode 设置日志模式
func (l *GormLogger) LogMode(level glog.LogLevel) glog.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

// Info 记录信息日志
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= glog.Info {
		l.log(ctx, glog.Info, msg, data...)
	}
}

// Warn 记录警告日志
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= glog.Warn {
		l.log(ctx, glog.Warn, msg, data...)
	}
}

// Error 记录错误日志
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= glog.Error {
		l.log(ctx, glog.Error, msg, data...)
	}
}

// Trace 记录SQL执行追踪
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= glog.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// add field
	fields := []logger.Field{
		logger.Str("sql", sql),
		logger.Int64("rows", rows),
		logger.Int64("duration_ms", elapsed.Milliseconds()),
	}

	switch {
	// 错误情况
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound):
		fields = append(fields, logger.Err(err))
		l.logger(ctx).Logc(2, logger.ERROR, "query error", fields...)
	// 慢查询
	case l.slowThreshold > 0 && elapsed > l.slowThreshold:
		l.logger(ctx).Logc(2, logger.WARN, "slow query", fields...)
	// 正常查询
	case l.logLevel >= glog.Info:
		l.logger(ctx).Logc(2, logger.INFO, "query success", fields...)
	}
}

// log 通用日志记录方法
func (l *GormLogger) log(ctx context.Context, level glog.LogLevel, msg string, data ...interface{}) {
	if len(data) > 0 {
		msg = fmt.Sprintf("%s | data: %v", msg, data)
	}

	switch level {
	case glog.Info:
		l.logger(ctx).Logc(2, logger.INFO, msg)
	case glog.Warn:
		l.logger(ctx).Logc(2, logger.WARN, msg)
	case glog.Error:
		l.logger(ctx).Logc(2, logger.ERROR, msg)
	default:
		l.logger(ctx).Logc(2, logger.INFO, msg)
	}
}

// logger 获取 logger，优先从 context 获取
func (l *GormLogger) logger(ctx context.Context) logger.Logger {
	if cl := logger.Ctx(ctx); cl != nil {
		return cl
	}
	return logger.Default()
}
