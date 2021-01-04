package helper

import "context"

type contextKey int

const (
	managedByLabelValueContextKey contextKey = iota
)

func ContextWithManagedByLabelValue(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, managedByLabelValueContextKey, value)
}

func ManagedByLabelValueFromContext(ctx context.Context) (value string, found bool) {
	value, found = ctx.Value(managedByLabelValueContextKey).(string)
	return
}
