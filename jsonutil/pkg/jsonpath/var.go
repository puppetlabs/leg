package jsonpath

import (
	"context"
	"strconv"

	"github.com/PaesslerAG/gval"
	"github.com/puppetlabs/leg/gvalutil/pkg/eval"
)

type variableWildcard struct{}

type variableRecursiveDescent struct{}

type variableRange struct {
	Min, Max, Step int
}

type variableChild struct {
	Key any
}

type VariableVisitor interface {
	VisitWildcard(ctx context.Context, parameter any, next func(context.Context, []PathValue) error) error
	VisitRecursiveDescent(ctx context.Context, parameter any, next func(context.Context, []PathValue) error) error
	VisitRange(ctx context.Context, parameter any, min, max, step int, next func(context.Context, []PathValue) error) error
	VisitChild(ctx context.Context, parameter any, key any, next func(context.Context, PathValue) error) error
}

type VariableVisitorFuncs struct {
	VisitWildcardFunc         func(ctx context.Context, parameter any, next func(context.Context, []PathValue) error) error
	VisitRecursiveDescentFunc func(ctx context.Context, parameter any, next func(context.Context, []PathValue) error) error
	VisitRangeFunc            func(ctx context.Context, parameter any, min, max, step int, next func(context.Context, []PathValue) error) error
	VisitChildFunc            func(ctx context.Context, parameter any, key any, next func(context.Context, PathValue) error) error
}

var _ VariableVisitor = VariableVisitorFuncs{}

func (vf VariableVisitorFuncs) VisitWildcard(c context.Context, v any, next func(context.Context, []PathValue) error) error {
	if vf.VisitWildcardFunc != nil {
		return vf.VisitWildcardFunc(c, v, next)
	}

	var items []PathValue

	appender := func(c context.Context, kv PathValue) error {
		items = append(items, kv)
		return nil
	}

	switch vt := v.(type) {
	case []any:
		for i := range vt {
			if err := vf.VisitChild(c, v, i, appender); err != nil {
				return &IndexParseError{Index: i, Cause: err}
			}
		}
	case map[string]any:
		for k := range vt {
			if err := vf.VisitChild(c, v, k, appender); err != nil {
				return &KeyParseError{Key: k, Cause: err}
			}
		}
	}

	return next(c, items)
}

func (vf VariableVisitorFuncs) VisitRecursiveDescent(c context.Context, v any, next func(context.Context, []PathValue) error) error {
	if vf.VisitRecursiveDescentFunc != nil {
		return vf.VisitRecursiveDescentFunc(c, v, next)
	}

	items := []PathValue{
		{Value: v},
	}

	var appender func(prefix ...string) func(context.Context, []PathValue) error
	appender = func(prefix ...string) func(context.Context, []PathValue) error {
		return func(c context.Context, vs []PathValue) error {
			for _, v := range vs {
				item := PathValue{
					Path:  append(append([]string{}, prefix...), v.Path...),
					Value: v.Value,
				}

				items = append(items, item)
				if err := vf.VisitWildcard(c, v.Value, appender(item.Path...)); err != nil {
					return &PathResolutionError{Path: v.Path, Cause: err}
				}
			}

			return nil
		}
	}
	if err := vf.VisitWildcard(c, v, appender()); err != nil {
		return err
	}

	return next(c, items)
}

func (vf VariableVisitorFuncs) VisitRange(c context.Context, v any, min, max, step int, next func(context.Context, []PathValue) error) error {
	if vf.VisitRangeFunc != nil {
		return vf.VisitRangeFunc(c, v, min, max, step, next)
	}

	var items []PathValue

	switch vt := v.(type) {
	case []any:
		n := len(vt)
		min = negmax(min, n)
		max = negmax(max, n)
	case map[string]any:
		// Ranging over a map is explicitly not supported.
		return next(c, items)
	default:
		// Otherwise we hope the variable selector can do its job.
	}

	if min > max {
		return nil
	}

	if step == 0 {
		step = 1
	}

	appender := func(c context.Context, kv PathValue) error {
		items = append(items, kv)
		return nil
	}

	if step > 0 {
		for i := min; i < max; i += step {
			if err := vf.VisitChild(c, v, i, appender); err != nil {
				return &IndexParseError{Index: i, Cause: err}
			}
		}
	} else {
		for i := max - 1; i >= min; i += step {
			if err := vf.VisitChild(c, v, i, appender); err != nil {
				return &IndexParseError{Index: i, Cause: err}
			}
		}
	}

	return next(c, items)
}

func negmax(n, max int) int {
	if n < 0 {
		n = max + n
		if n < 0 {
			n = 0
		}
	} else if n > max {
		return max
	}
	return n
}

func (vf VariableVisitorFuncs) VisitChild(c context.Context, v, key any, next func(context.Context, PathValue) error) error {
	if vf.VisitChildFunc != nil {
		return vf.VisitChildFunc(c, v, key, next)
	}

	switch vt := v.(type) {
	case []any:
		i, r, err := eval.SelectIndex(vt, key)
		if err != nil {
			return err
		}

		return next(c, PathValue{
			Path:  []string{strconv.Itoa(i)},
			Value: r,
		})
	case map[string]any:
		k, r, err := eval.SelectKey(vt, key)
		if err != nil {
			return err
		}

		return next(c, PathValue{
			Path:  []string{k},
			Value: r,
		})
	default:
		return &eval.UnsupportedValueTypeError{Value: vt}
	}
}

func VariableSelector(visitor VariableVisitor) func(path gval.Evaluables) gval.Evaluable {
	return func(path gval.Evaluables) gval.Evaluable {
		return func(c context.Context, v any) (r any, err error) {
			var next func(c context.Context, rest gval.Evaluables, v any) (values, bool, error)
			next = func(c context.Context, rest gval.Evaluables, v any) (values, bool, error) {
				if len(rest) == 0 {
					return nil, false, nil
				}

				t, err := rest[0](c, v)
				if err != nil {
					return nil, false, err
				}

				var r values
				combine := func(c context.Context, vs values) error {
					vs, err := vs.flatMap(func(v value) (values, error) {
						r, ok, err := next(c, rest[1:], v.value)
						if err != nil {
							return nil, err
						} else if !ok {
							return v, nil
						}

						return v.prefix(r), nil
					})
					if err != nil {
						return err
					}

					vs.concat(&r)
					return nil
				}

				switch tt := t.(type) {
				case variableWildcard:
					err = visitor.VisitWildcard(c, v, func(c context.Context, pvs []PathValue) error {
						return combine(c, pathValueSlice(pvs))
					})
				case variableRecursiveDescent:
					err = visitor.VisitRecursiveDescent(c, v, func(c context.Context, pvs []PathValue) error {
						return combine(c, pathValueSlice(pvs))
					})
				case variableRange:
					err = visitor.VisitRange(c, v, tt.Min, tt.Max, tt.Step, func(c context.Context, pvs []PathValue) error {
						return combine(c, pathValueSlice(pvs))
					})
				case variableChild:
					err = visitor.VisitChild(c, v, tt.Key, func(c context.Context, pv PathValue) error {
						return combine(c, value{wildcards: [][]string{pv.Path}, value: pv.Value})
					})
				default:
					err = &UnknownVariableTypeError{Variable: t}
				}
				return r, err == nil, err
			}
			r, _, err = next(c, path, v)
			return
		}
	}
}

func pathValueSlice(pvs []PathValue) valueSlice {
	vs := make(valueSlice, len(pvs))
	for i, pv := range pvs {
		vs[i] = value{wildcards: [][]string{pv.Path}, value: pv.Value}
	}
	return vs
}

func ChildVariableSelector(fn func(ctx context.Context, parameter any, key any, next func(context.Context, PathValue) error) error) func(path gval.Evaluables) gval.Evaluable {
	return VariableSelector(VariableVisitorFuncs{
		VisitChildFunc: fn,
	})
}

func DefaultVariableVisitor() VariableVisitor {
	return VariableVisitorFuncs{}
}

func DefaultVariableSelector() func(path gval.Evaluables) gval.Evaluable {
	return VariableSelector(DefaultVariableVisitor())
}
