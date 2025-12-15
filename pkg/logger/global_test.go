package logger

import (
	"context"
	"errors"
	"testing"
)

func TestCtxInfo(t *testing.T) {
	Init()

	t.Run("CtxLog", func(t *testing.T) {
		ctx := context.Background()
		Info(ctx, "test")

		ctx = CtxWithField(ctx, Str("requestId", "23123"))

		Warn(ctx, "test", Int("age", 10))
		Errw(ctx, errors.New("test error"), Int("age", 18))

		log := Ctx(ctx)
		ctx = CtxWithLogger(ctx, log)
		Error(ctx, "test")

		log.Log(INFO, "test", Str("add", "new add"))
		log.Logc(-1, WARN, "test", Str("add", "new add"))
	})
}
