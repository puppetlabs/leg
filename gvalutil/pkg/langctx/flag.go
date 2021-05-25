package langctx

import (
	"context"
	"sync/atomic"
)

type contextValue uint64

var contextValueCounter uint64

type BoolFlag struct {
	id contextValue
}

func (bf *BoolFlag) With(ctx context.Context) (nctx context.Context, v *bool) {
	v = new(bool)
	nctx = bf.WithVar(ctx, v)
	return
}

func (bf *BoolFlag) WithVar(ctx context.Context, v *bool) context.Context {
	return context.WithValue(ctx, bf.id, v)
}

func (bf *BoolFlag) Set(ctx context.Context) bool {
	v, ok := ctx.Value(bf.id).(*bool)
	if !ok || v == nil {
		return false
	}

	*v = true
	return true
}

func (bf *BoolFlag) MustSet(ctx context.Context) {
	if !bf.Set(ctx) {
		panic("flag not in context")
	}
}

func NewBoolFlag() *BoolFlag {
	return &BoolFlag{
		id: contextValue(atomic.AddUint64(&contextValueCounter, 1)),
	}
}
