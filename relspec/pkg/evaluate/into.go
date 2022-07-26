package evaluate

import (
	"context"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

func mapstructureHookFunc[T ref.References[T]](ctx context.Context, ev Evaluator[T], md *Metadata[T]) mapstructure.DecodeHookFunc {
	return func(from reflect.Type, to reflect.Type, data any) (any, error) {
		depth := DepthFull

		// Copy so we can potentially use the zero value below.
		check := to
		for check.Kind() == reflect.Ptr {
			check = check.Elem()
		}

		if check.Kind() == reflect.Struct {
			// We only evaluate one level of nesting for structs, because their
			// children will get correctly traversed once the data exists.
			depth = 1
		}

		r, err := ev.Evaluate(ctx, data, depth)
		if err != nil {
			return nil, err
		}

		// TODO: Is there a way we can possibly preserve information about the
		// evaluator to use when recursing into struct fields?
		md.MergeReferences(r.Metadata)

		if !r.OK() {
			// We return the zero value of the type to eliminate confusion.
			return reflect.Zero(to).Interface(), nil
		}

		return r.Value, nil
	}
}

func EvaluateInto[T ref.References[T]](ctx context.Context, ev Evaluator[T], from, to any) (*Metadata[T], error) {
	md := EmptyMetadata[T]()

	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructureHookFunc(ctx, ev, md),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339Nano),
		),
		ZeroFields: true,
		Result:     to,
		TagName:    "spec",
	})
	if err != nil {
		return nil, err
	}

	return md, d.Decode(from)
}
