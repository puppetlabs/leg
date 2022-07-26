package pathlang

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"regexp"

	"github.com/PaesslerAG/gval"
	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/query"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
	"golang.org/x/exp/constraints"
)

type convFunc[T any] func(x any) (T, bool)

func convCast[T any](x any) (v T, ok bool) {
	v, ok = x.(T)
	return
}

func convString(x any) (string, bool) {
	if v, ok := convCast[string](x); ok {
		return v, true
	}
	if x == nil {
		return "", false
	}
	return fmt.Sprintf("%v", x), true
}

func convNumber(x any) (float64, bool) {
	f := gval.Evaluable(func(_ context.Context, _ any) (any, error) {
		return x, nil
	})
	v, err := f.EvalFloat64(context.Background(), nil)
	return v, err == nil
}

func convBool(x any) (value, converted bool) {
	f := gval.Evaluable(func(_ context.Context, _ any) (any, error) {
		return x, nil
	})
	v, err := f.EvalBool(context.Background(), nil)
	return v, err == nil
}

type op1Func[T ref.References[T]] func(a *evaluate.Result[T]) (value any, ok bool, err error)
type apply1Func[T any] func(a T) (any, error)

func eval1[T ref.References[T]](r *evaluate.Result[T], fn op1Func[T]) (any, error) {
	if v, ok, err := fn(r); err != nil || ok {
		return v, err
	}
	return nil, fmt.Errorf("operand of type %T is not supported", r.Value)
}

func ifConvertibleToType1[R ref.References[R], T any](conv convFunc[T], fn apply1Func[T]) op1Func[R] {
	return func(a *evaluate.Result[R]) (any, bool, error) {
		v1, ok := conv(a.Value)
		if !ok {
			return nil, false, nil
		}

		v, err := fn(v1)
		return v, true, err
	}
}

func passIfConvertibleToType1[R ref.References[R], T any](conv convFunc[T]) op1Func[R] {
	return ifConvertibleToType1[R](conv, func(a T) (any, error) { return a, nil })
}

type op2Func[T ref.References[T]] func(a, b *evaluate.Result[T]) (any, bool, error)
type apply2Func[T any] func(a, b T) (any, error)

func plus[T constraints.Ordered](a, b T) (any, error)               { return a + b, nil }
func equals[T comparable](a, b T) (any, error)                      { return a == b, nil }
func notEquals[T comparable](a, b T) (any, error)                   { return a != b, nil }
func greaterThan[T constraints.Ordered](a, b T) (any, error)        { return a > b, nil }
func greaterThanOrEqual[T constraints.Ordered](a, b T) (any, error) { return a >= b, nil }
func lessThan[T constraints.Ordered](a, b T) (any, error)           { return a < b, nil }
func lessThanOrEqual[T constraints.Ordered](a, b T) (any, error)    { return a <= b, nil }

func eval2[T ref.References[T]](a, b *evaluate.Result[T], fn op2Func[T]) (any, error) {
	if v, ok, err := fn(a, b); err != nil || ok {
		return v, err
	}
	return nil, fmt.Errorf("operands of types %T and %T are not supported", a.Value, b.Value)
}

func chain2[T ref.References[T]](fns ...op2Func[T]) op2Func[T] {
	return func(a, b *evaluate.Result[T]) (any, bool, error) {
		for _, fn := range fns {
			v, ok, err := fn(a, b)
			if err != nil || ok {
				return v, ok, err
			}
		}
		return nil, false, nil
	}
}

func ifConvertibleToType2[R ref.References[R], T any](conv convFunc[T], fn apply2Func[T]) op2Func[R] {
	return func(a, b *evaluate.Result[R]) (any, bool, error) {
		v1, ok1 := conv(a.Value)
		v2, ok2 := conv(b.Value)
		if !ok1 || !ok2 {
			return nil, false, nil
		}

		v, err := fn(v1, v2)
		return v, true, err
	}
}

func ifOfType2[R ref.References[R], T any](fn apply2Func[T]) op2Func[R] {
	return ifConvertibleToType2[R](convCast[T], fn)
}

func interleaveIfConvertibleToType2For2[R ref.References[R], T, U any](conv1 convFunc[T], fn1 apply2Func[T], conv2 convFunc[U], fn2 apply2Func[U]) op2Func[R] {
	return chain2(
		ifOfType2[R](fn1),
		ifOfType2[R](fn2),
		ifConvertibleToType2[R](conv1, fn1),
		ifConvertibleToType2[R](conv2, fn2),
	)
}

type unaryFunc[T ref.References[T]] func(r *evaluate.Result[T]) (any, error)

func unaryEval[T ref.References[T]](fn op1Func[T]) unaryFunc[T] {
	return func(r *evaluate.Result[T]) (any, error) {
		return eval1(r, fn)
	}
}

type binaryFunc[T ref.References[T]] func(ra, rb *evaluate.Result[T]) (any, error)

func binaryEval[T ref.References[T]](fn op2Func[T]) binaryFunc[T] {
	return func(ra, rb *evaluate.Result[T]) (any, error) {
		return eval2(ra, rb, fn)
	}
}

type condFunc[T ref.References[T]] func(r *evaluate.Result[T]) (any, bool, error)

func condEval[T ref.References[T]](fn op1Func[T], check func(v any) bool) condFunc[T] {
	return func(r *evaluate.Result[T]) (any, bool, error) {
		v, ok, err := fn(r)
		if err != nil || !ok {
			return nil, false, err
		}

		return v, check(v), nil
	}
}

func condNever[T ref.References[T]](r *evaluate.Result[T]) (value any, ok bool, err error) {
	return nil, false, nil
}

var _ condFunc[ref.EmptyReferences] = condNever[ref.EmptyReferences]

type operatorDefs[T ref.References[T]] struct {
	evaluator       evaluate.Evaluator[T]
	eagerEvaluation bool
}

func (od *operatorDefs[T]) InfixShortCircuit(name string, lfn condFunc[T], fn binaryFunc[T]) gval.Language {
	return gval.InfixEvalOperator(name, func(a, b gval.Evaluable) (gval.Evaluable, error) {
		if a.IsConst() && b.IsConst() {
			va, err := a(context.Background(), nil)
			if err != nil {
				return nil, err
			}

			vb, err := b(context.Background(), nil)
			if err != nil {
				return nil, err
			}

			v, err := fn(evaluate.StaticResult[T](va), evaluate.StaticResult[T](vb))
			if err != nil {
				return nil, err
			}

			return func(ctx context.Context, parameter any) (any, error) {
				return v, nil
			}, nil
		}

		return func(ctx context.Context, parameter any) (v any, err error) {
			ctx, md := evaluate.ContextWithEmptyMetadata[T](ctx)
			md.SetAnnotation(InfixOperationAnnotation(name))

			ra, err := query.EvaluateGval(ctx, evaluate.CoalesceEvaluators(md.Evaluator, od.evaluator), a, parameter)
			if err != nil {
				return nil, err
			}

			md.AddStep(ra, OperandAnnotation(1))

			var ok bool
			if md.OK() {
				v, ok, err = lfn(ra)
				if err != nil {
					return nil, err
				}
			}

			if !ok || od.eagerEvaluation {
				rb, err := query.EvaluateGval(ctx, evaluate.CoalesceEvaluators(md.Evaluator, od.evaluator), b, parameter)
				if err != nil {
					return nil, err
				}

				rb.References.SetUsed(!ok)
				md.AddStep(rb, OperandAnnotation(2))

				if !md.OK() {
					return evaluate.StaticExpandable(evaluate.ContextualizedResult(md)), nil
				}

				if !ok {
					v, err = fn(ra, rb)
					if err != nil {
						return nil, err
					}
				}
			}

			return evaluate.StaticExpandable(evaluate.NewResult(md, v)), nil
		}, nil
	})
}

func (od *operatorDefs[T]) Infix(name string, fn binaryFunc[T]) gval.Language {
	return od.InfixShortCircuit(name, condNever[T], fn)
}

func (od *operatorDefs[T]) Prefix(name rune, fn unaryFunc[T]) gval.Language {
	return gval.PrefixExtension(name, func(ctx context.Context, p *gval.Parser) (gval.Evaluable, error) {
		eval, err := p.ParseNextExpression(ctx)
		if err != nil {
			return nil, err
		}

		if eval.IsConst() {
			v, err := eval(context.Background(), nil)
			if err != nil {
				return nil, err
			}

			v, err = fn(evaluate.StaticResult[T](v))
			if err != nil {
				return nil, err
			}

			return p.Const(v), nil
		}

		return func(ctx context.Context, parameter any) (any, error) {
			ctx, md := evaluate.ContextWithEmptyMetadata[T](ctx)
			md.SetAnnotation(PrefixOperationAnnotation(string(name)))

			r, err := query.EvaluateGval(ctx, evaluate.CoalesceEvaluators(md.Evaluator, od.evaluator), eval, parameter)
			if err != nil {
				return nil, err
			}

			md.AddStep(r, OperandAnnotation(1))

			if !md.OK() {
				return evaluate.StaticExpandable(evaluate.ContextualizedResult(md)), nil
			}

			v, err := fn(r)
			if err != nil {
				return nil, err
			}

			return evaluate.StaticExpandable(evaluate.NewResult(md, v)), nil
		}, nil
	})
}

func operators[T ref.References[T]](ev evaluate.Evaluator[T], o *Options[T]) gval.Language {
	defs := &operatorDefs[T]{
		evaluator:       ev,
		eagerEvaluation: o.EagerEvaluation,
	}

	lang := gval.NewLanguage(
		defs.Prefix('!', unaryEval(ifConvertibleToType1[T](convBool, func(a bool) (any, error) { return !a, nil }))),
		defs.InfixShortCircuit(
			"&&",
			condEval(passIfConvertibleToType1[T](convBool), func(v any) bool {
				return v == false
			}),
			binaryEval(ifConvertibleToType2[T](convBool, func(a, b bool) (any, error) { return a && b, nil })),
		),
		defs.InfixShortCircuit(
			"||",
			condEval(passIfConvertibleToType1[T](convBool), func(v any) bool {
				return v == true
			}),
			binaryEval(ifConvertibleToType2[T](convBool, func(a, b bool) (any, error) { return a || b, nil })),
		),

		defs.Infix("==", binaryEval(chain2(
			interleaveIfConvertibleToType2For2[T](convNumber, equals[float64], convBool, equals[bool]),
			func(a, b *evaluate.Result[T]) (any, bool, error) {
				return reflect.DeepEqual(a.Value, b.Value), true, nil
			},
		))),
		defs.Infix("!=", binaryEval(chain2(
			interleaveIfConvertibleToType2For2[T](convNumber, notEquals[float64], convBool, notEquals[bool]),
			func(a, b *evaluate.Result[T]) (any, bool, error) {
				return !reflect.DeepEqual(a.Value, b.Value), true, nil
			},
		))),
		defs.Infix(">", binaryEval(
			interleaveIfConvertibleToType2For2[T](convNumber, greaterThan[float64], convString, greaterThan[string]),
		)),
		defs.Infix(">=", binaryEval(
			interleaveIfConvertibleToType2For2[T](convNumber, greaterThanOrEqual[float64], convString, greaterThanOrEqual[string]),
		)),
		defs.Infix("<", binaryEval(
			interleaveIfConvertibleToType2For2[T](convNumber, lessThan[float64], convString, lessThan[string]),
		)),
		defs.Infix("<=", binaryEval(
			interleaveIfConvertibleToType2For2[T](convNumber, lessThanOrEqual[float64], convString, lessThanOrEqual[string]),
		)),

		defs.Infix("+", binaryEval(
			interleaveIfConvertibleToType2For2[T](convNumber, plus[float64], convString, plus[string]),
		)),

		defs.Infix("=~", binaryEval(
			ifConvertibleToType2[T](convString, func(a, b string) (any, error) { return regexp.MatchString(b, a) }),
		)),
		defs.Infix("!~", binaryEval(
			ifConvertibleToType2[T](convString, func(a, b string) (any, error) {
				m, err := regexp.MatchString(b, a)
				if err != nil {
					return nil, err
				}
				return !m, nil
			}),
		)),

		defs.Prefix('-', unaryEval(ifConvertibleToType1[T](convNumber, func(a float64) (any, error) { return -a, nil }))),
		defs.Infix("-", binaryEval(ifConvertibleToType2[T](convNumber, func(a, b float64) (any, error) { return a - b, nil }))),
		defs.Infix("*", binaryEval(ifConvertibleToType2[T](convNumber, func(a, b float64) (any, error) { return a * b, nil }))),
		defs.Infix("/", binaryEval(ifConvertibleToType2[T](convNumber, func(a, b float64) (any, error) { return a / b, nil }))),
		defs.Infix("%", binaryEval(ifConvertibleToType2[T](convNumber, func(a, b float64) (any, error) { return math.Mod(a, b), nil }))),
		defs.Infix("**", binaryEval(ifConvertibleToType2[T](convNumber, func(a, b float64) (any, error) { return math.Pow(a, b), nil }))),

		defs.Prefix('~', unaryEval(ifConvertibleToType1[T](convNumber, func(a float64) (any, error) { return float64(^int64(a)), nil }))),
		defs.Infix("^", binaryEval(
			ifConvertibleToType2[T](convNumber, func(a, b float64) (any, error) { return float64(int64(a) ^ int64(b)), nil }),
		)),
		defs.Infix("&", binaryEval(
			ifConvertibleToType2[T](convNumber, func(a, b float64) (any, error) { return float64(int64(a) & int64(b)), nil }),
		)),
		defs.Infix("|", binaryEval(
			ifConvertibleToType2[T](convNumber, func(a, b float64) (any, error) { return float64(int64(a) | int64(b)), nil }),
		)),
		defs.Infix("<<", binaryEval(
			ifConvertibleToType2[T](convNumber, func(a, b float64) (any, error) { return float64(int64(a) << uint64(b)), nil }),
		)),
		defs.Infix(">>", binaryEval(
			ifConvertibleToType2[T](convNumber, func(a, b float64) (any, error) { return float64(int64(a) >> uint64(b)), nil }),
		)),
	)

	return lang
}
