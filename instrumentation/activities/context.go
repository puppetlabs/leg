package activities

import "context"

type contextKey int

const (
	reporterKey contextKey = iota
)

func NewContext(ctx context.Context, reporter *Reporter) context.Context {
	return context.WithValue(ctx, reporterKey, reporter)
}

func FromContext(ctx context.Context) *Reporter {
	obj := ctx.Value(reporterKey)

	if r, ok := obj.(*Reporter); ok {
		return r
	}

	return nil
}
