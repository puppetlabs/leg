package eventctx

import (
	"context"

	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type contextKey int

const (
	eventRecorderContextKey contextKey = iota
)

var (
	stubRecorder = &record.FakeRecorder{}
)

func WithEventRecorder(ctx context.Context, mgr manager.Manager, name string) context.Context {
	return context.WithValue(ctx, eventRecorderContextKey, mgr.GetEventRecorderFor(name))
}

func EventRecorder(ctx context.Context) record.EventRecorder {
	r, ok := ctx.Value(eventRecorderContextKey).(record.EventRecorder)
	if !ok {
		r = stubRecorder
	}

	return r
}
