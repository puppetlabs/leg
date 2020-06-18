package activities

import "context"

const contextKey = "horsehead.activities.reporter"

func NewContext(ctx context.Context, reporter *Reporter) context.Context {
	return context.WithValue(ctx, contextKey, reporter)
}

func FromContext(ctx context.Context) *Reporter {
	obj := ctx.Value(contextKey)

	if r, ok := obj.(*Reporter); ok {
		return r
	}

	return nil
}
