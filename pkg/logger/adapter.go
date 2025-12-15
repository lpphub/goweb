package logger

import "context"

type ctxAdapter func(context.Context) context.Context

var ctxAdapters []ctxAdapter

func RegisterCtxAdapter(adapter ctxAdapter) {
	ctxAdapters = append(ctxAdapters, adapter)
}
