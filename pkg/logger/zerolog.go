package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type zeroLogger struct{ core zerolog.Logger }

func newZeroLogger(cfg *config) Logger {
	// default config
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.CallerMarshalFunc = callerShortFunc

	output := cfg.output
	if output == nil {
		output = os.Stdout
	}

	l := zerolog.New(output).With().Timestamp().Logger().Level(zerolog.Level(cfg.level))
	return &zeroLogger{core: l}
}

func (z *zeroLogger) Log(level Level, msg string, fields ...Field) {
	z.Logc(1, level, msg, fields...)
}

func (z *zeroLogger) Logc(callerSkip int, level Level, msg string, fields ...Field) {
	e := z.event(level)

	// skip frame
	e = e.CallerSkipFrame(callerSkip + 2)

	for _, f := range fields {
		e = z.addField(e, f)
	}
	e.Msg(msg)
}

func (z *zeroLogger) With(fields ...Field) Logger {
	ctx := z.core.With()
	for _, f := range fields {
		ctx = z.withField(ctx, f)
	}
	return &zeroLogger{core: ctx.Logger()}
}

func (z *zeroLogger) WithCaller(skip int) Logger {
	return &zeroLogger{core: z.core.With().CallerWithSkipFrameCount(skip).Logger()}
}

func (z *zeroLogger) event(level Level) *zerolog.Event {
	switch level {
	case DEBUG:
		return z.core.Debug()
	case WARN:
		return z.core.Warn()
	case ERROR:
		return z.core.Error()
	case FATAL:
		return z.core.Fatal()
	default:
		return z.core.Info()
	}
}

func (z *zeroLogger) addField(e *zerolog.Event, f Field) *zerolog.Event {
	switch v := f.Value.(type) {
	case string:
		return e.Str(f.Key, v)
	case int:
		return e.Int(f.Key, v)
	case int64:
		return e.Int64(f.Key, v)
	case float64:
		return e.Float64(f.Key, v)
	case bool:
		return e.Bool(f.Key, v)
	case error:
		return e.AnErr(f.Key, v)
	default:
		return e.Interface(f.Key, v)
	}
}

func (z *zeroLogger) withField(c zerolog.Context, f Field) zerolog.Context {
	switch v := f.Value.(type) {
	case string:
		return c.Str(f.Key, v)
	case int:
		return c.Int(f.Key, v)
	case int64:
		return c.Int64(f.Key, v)
	case float64:
		return c.Float64(f.Key, v)
	case bool:
		return c.Bool(f.Key, v)
	case error:
		return c.AnErr(f.Key, v)
	default:
		return c.Interface(f.Key, v)
	}
}

func callerShortFunc(_ uintptr, file string, line int) string {
	file = filepath.ToSlash(file)
	parts := strings.Split(file, "/")
	if len(parts) > 2 {
		file = strings.Join(parts[len(parts)-2:], "/")
	}
	return fmt.Sprintf("%s:%d", file, line)
}
