package evaluate

import (
	"context"
	"reflect"
	"strconv"
	"time"

	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

const DepthFull = -1

func IsFullDepth(depth int) bool {
	return depth < 0
}

type Evaluator[T ref.References[T]] interface {
	Evaluate(ctx context.Context, tree any, depth int) (*Result[T], error)
}

type Visitor[T ref.References[T]] interface {
	VisitExpandable(ctx context.Context, ex Expandable[T], depth int, next Evaluator[T]) (*Result[T], error)
	VisitSlice(ctx context.Context, s []any, depth int, next Evaluator[T]) (*Result[T], error)
	VisitMap(ctx context.Context, m map[string]any, depth int, next Evaluator[T]) (*Result[T], error)
	VisitString(ctx context.Context, s string, depth int, next Evaluator[T]) (*Result[T], error)
}

type VisitorFuncs[T ref.References[T]] struct {
	VisitExpandableFunc func(ctx context.Context, ex Expandable[T], depth int, next Evaluator[T]) (*Result[T], error)
	VisitSliceFunc      func(ctx context.Context, s []any, depth int, next Evaluator[T]) (*Result[T], error)
	VisitMapFunc        func(ctx context.Context, m map[string]any, depth int, next Evaluator[T]) (*Result[T], error)
	VisitStringFunc     func(ctx context.Context, s string, depth int, next Evaluator[T]) (*Result[T], error)
}

func (vfs *VisitorFuncs[T]) VisitExpandable(ctx context.Context, ex Expandable[T], depth int, next Evaluator[T]) (*Result[T], error) {
	if vfs.VisitExpandableFunc != nil {
		return vfs.VisitExpandableFunc(ctx, ex, depth, next)
	}

	r := EmptyResult[T]()
	r.SetAnnotation(ExpansionAnnotation)

	pr, err := ex.Expand(ctx, depth)
	if err != nil {
		return nil, err
	}

	r.AddStep(pr)
	r.MergeEvaluator(pr.Metadata)
	r.SetValue(pr.Value)

	if !r.OK() {
		return r, nil
	}

	// It's possible for an expandable to return another expandable, in which
	// case we'll continue to evaluate at this depth.
	nr, err := CoalesceEvaluators(r.Evaluator, next).Evaluate(ctx, pr.Value, depth)
	if err != nil {
		return nil, err
	}

	r.AddStep(nr)
	r.MergeEvaluator(nr.Metadata)

	if r.OK() {
		r.SetValue(nr.Value)
	}

	return r, nil
}

func (vfs *VisitorFuncs[T]) VisitSlice(ctx context.Context, s []any, depth int, next Evaluator[T]) (*Result[T], error) {
	if vfs.VisitSliceFunc != nil {
		return vfs.VisitSliceFunc(ctx, s, depth, next)
	}

	r := EmptyResult[T]()
	r.SetAnnotation(ArrayAnnotation)

	if depth == 1 {
		r.SetValue(s)
		return r, nil
	}

	l := make([]any, len(s))
	for i, v := range s {
		nv, err := next.Evaluate(ctx, v, depth-1)
		if err != nil {
			return nil, &PathEvaluationError{
				Path:  strconv.Itoa(i),
				Cause: err,
			}
		}

		r.AddStep(nv, ArrayIndexAnnotation(i))
		l[i] = nv.Value
	}

	r.SetValue(l)
	return r, nil
}

func (vfs *VisitorFuncs[T]) VisitMap(ctx context.Context, m map[string]any, depth int, next Evaluator[T]) (*Result[T], error) {
	if vfs.VisitMapFunc != nil {
		return vfs.VisitMapFunc(ctx, m, depth, next)
	}

	r := EmptyResult[T]()
	r.SetAnnotation(ObjectAnnotation)

	if depth == 1 {
		r.SetValue(m)
		return r, nil
	}

	rm := make(map[string]any, len(m))
	for k, v := range m {
		nv, err := next.Evaluate(ctx, v, depth-1)
		if err != nil {
			return nil, &PathEvaluationError{Path: k, Cause: err}
		}

		r.AddStep(nv, ObjectKeyAnnotation(k))
		rm[k] = nv.Value
	}

	r.SetValue(rm)
	return r, nil
}

func (vfs *VisitorFuncs[T]) VisitString(ctx context.Context, s string, depth int, next Evaluator[T]) (*Result[T], error) {
	if vfs.VisitStringFunc != nil {
		return vfs.VisitStringFunc(ctx, s, depth, next)
	}

	return StaticResult[T](s), nil
}

func DefaultVisitor[T ref.References[T]]() Visitor[T] {
	return &VisitorFuncs[T]{}
}

type visitorEvaluator[T ref.References[T]] struct {
	visitor Visitor[T]
}

func (ve *visitorEvaluator[T]) evaluate(ctx context.Context, tree any, depth int) (*Result[T], error) {
	if depth == 0 {
		return StaticResult[T](tree), nil
	}

	switch vt := tree.(type) {
	case Expandable[T]:
		return ve.visitor.VisitExpandable(ctx, vt, depth, ve)
	case []any:
		return ve.visitor.VisitSlice(ctx, vt, depth, ve)
	case map[string]any:
		return ve.visitor.VisitMap(ctx, vt, depth, ve)
	case string:
		return ve.visitor.VisitString(ctx, vt, depth, ve)
	default:
		return StaticResult[T](tree), nil
	}
}

func (ve *visitorEvaluator[T]) Evaluate(ctx context.Context, tree any, depth int) (*Result[T], error) {
	candidate, err := ve.evaluate(ctx, tree, depth)
	if err != nil {
		return nil, err
	}

	if depth == 0 {
		return candidate, nil
	}

	switch candidate.Value.(type) {
	// Valid JSON types per https://golang.org/pkg/encoding/json/:
	case bool, float64, string, []any, map[string]any, nil:
		return candidate, nil
	// We support a set of additional YAML scalar(-ish) types decoded by
	// gopkg.in/yaml.v3.
	case []byte, int, int64, uint, uint64, time.Time:
		return candidate, nil
	default:
		return nil, &UnsupportedValueError{Type: reflect.TypeOf(candidate.Value)}
	}
}

func NewEvaluator[T ref.References[T]](visitor Visitor[T]) Evaluator[T] {
	return &visitorEvaluator[T]{visitor: visitor}
}

func DefaultEvaluator[T ref.References[T]]() Evaluator[T] {
	return NewEvaluator(DefaultVisitor[T]())
}

func CoalesceEvaluators[T ref.References[T]](evs ...Evaluator[T]) Evaluator[T] {
	for _, ev := range evs {
		if ev != nil {
			return ev
		}
	}
	return DefaultEvaluator[T]()
}

func EvaluateAll[T ref.References[T]](ctx context.Context, ev Evaluator[T], tree any) (*Result[T], error) {
	return ev.Evaluate(ctx, tree, DepthFull)
}
