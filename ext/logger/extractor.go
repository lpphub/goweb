package logger

import "context"

type ctxExtractor func(context.Context) context.Context

var ctxExtractors []ctxExtractor

func RegisterCtxExtractor(extractor ctxExtractor) {
	ctxExtractors = append(ctxExtractors, extractor)
}
