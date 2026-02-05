package logx

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/lpphub/goweb/pkg/logging"
	"github.com/redis/go-redis/v9"
)

// RedisLogger 自定义Redis客户端日志记录器
type RedisLogger struct {
	slowThreshold time.Duration
	logger        logging.Logger
}

// NewRedisLogger 创建新的Redis日志记录器
func NewRedisLogger() *RedisLogger {
	return &RedisLogger{
		slowThreshold: 100 * time.Millisecond,
		logger:        logging.L().WithCaller(5),
	}
}

func (l *RedisLogger) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		start := time.Now()
		conn, err := next(ctx, network, addr)
		elapsed := time.Since(start)

		fields := map[string]interface{}{
			"addr":        addr,
			"duration_ms": l.fmtDuration(elapsed),
		}

		if err != nil {
			l.logger.Error(ctx).Fields(fields).Err(err).Msg("redis connected failed")
		} else {
			l.logger.Info(ctx).Fields(fields).Msg("redis connected")
		}
		return conn, err
	}
}

// ProcessHook 实现命令处理钩子
func (l *RedisLogger) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()
		err := next(ctx, cmd)
		elapsed := time.Since(start)

		// 添加字段
		fields := map[string]interface{}{
			"cmd":         l.buildCmd(cmd),
			"duration_ms": l.fmtDuration(elapsed),
		}

		switch {
		case err != nil && !errors.Is(err, redis.Nil):
			l.logger.Error(ctx).Fields(fields).Err(err).Msg("redis error")
		case l.slowThreshold > 0 && elapsed > l.slowThreshold:
			l.logger.Warn(ctx).Fields(fields).Msg("redis slow")
		default:
			l.logger.Info(ctx).Fields(fields).Msg("redis success")
		}
		return err
	}
}

// ProcessPipelineHook 实现管道处理钩子
func (l *RedisLogger) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		start := time.Now()
		err := next(ctx, cmds)
		elapsed := time.Since(start)

		// 记录管道执行的整体信息
		fields := map[string]interface{}{
			"cmd":         l.buildPipelineCmd(cmds),
			"duration_ms": l.fmtDuration(elapsed),
		}

		switch {
		case err != nil && !errors.Is(err, redis.Nil):
			l.logger.Error(ctx).Fields(fields).Err(err).Msg("redis pipeline error")
		case l.slowThreshold > 0 && elapsed > l.slowThreshold:
			l.logger.Warn(ctx).Fields(fields).Msg("redis pipeline slow")
		default:
			l.logger.Info(ctx).Fields(fields).Msg("redis pipeline success")
		}
		return err
	}
}

// buildCommandString 构建命令字符串（隐藏敏感信息）
func (l *RedisLogger) buildCmd(cmd redis.Cmder) string {
	args := cmd.Args()
	if len(args) == 0 {
		return cmd.Name()
	}

	// 转换参数为字符串
	var argStrs []string
	for _, arg := range args {
		argStrs = append(argStrs, fmt.Sprintf("%v", arg))
	}

	// 隐藏敏感命令的参数
	cmdName := strings.ToUpper(cmd.Name())
	if l.isSensitiveCommand(cmdName) {
		if len(argStrs) > 1 {
			// 保留命令名，隐藏参数
			return fmt.Sprintf("%s ****", argStrs[0])
		}
	}

	return strings.Join(argStrs, " ")
}

func (l *RedisLogger) buildPipelineCmd(cmds []redis.Cmder) string {
	if len(cmds) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, cmd := range cmds {
		if i > 0 {
			sb.WriteString("; ")
		}
		if i >= 5 {
			sb.WriteString(fmt.Sprintf("... (%d more)", len(cmds)-5))
			break
		}
		sb.WriteString(l.buildCmd(cmd))
	}
	return sb.String()
}

func (l *RedisLogger) fmtDuration(d time.Duration) string {
	return fmt.Sprintf("%.3f", float64(d.Nanoseconds())/1e6)
}

// isSensitiveCommand 判断是否为敏感命令
func (l *RedisLogger) isSensitiveCommand(cmdName string) bool {
	sensitiveCommands := map[string]bool{
		"AUTH":    true,
		"SET":     false, // 可能包含敏感值
		"GET":     false,
		"HSET":    false, // 可能包含敏感值
		"HGET":    false,
		"LPUSH":   false, // 可能包含敏感值
		"RPUSH":   false, // 可能包含敏感值
		"LSET":    false, // 可能包含敏感值
		"CONFIG":  true,  // 可能包含敏感配置
		"DEBUG":   true,
		"EVAL":    true, // 脚本可能包含敏感信息
		"SCRIPT":  true,
		"MIGRATE": true,
		"RESTORE": true,
	}
	return sensitiveCommands[cmdName]
}
