package fn

import (
	"context"

	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
)

func evaluateArg[T ref.References[T]](ctx context.Context, ev evaluate.Evaluator[T], arg, value any, depth int) (*evaluate.Result[T], error) {
	md := evaluate.MetadataFromContext[T](ctx)

	r, err := evaluate.CoalesceEvaluators(md.Evaluator, ev).Evaluate(ctx, value, depth)
	if err != nil {
		return nil, err
	}

	return r, nil
}

type Arg[T ref.References[T]] interface {
	evaluate.Expandable[T]
	Get() (any, bool)
	Annotation() evaluate.Annotation
	WrapError(err error) error
}

type ComparableArg[T ref.References[T]] interface {
	comparable
	Arg[T]
}

func ArgValue[T ref.References[T], A Arg[T]](arg A) any {
	v, _ := arg.Get()
	return v
}

func ArgPresent[T ref.References[T], A Arg[T]](arg A) bool {
	_, found := arg.Get()
	return found
}

type Args[T ref.References[T], A Arg[T], V any] interface {
	InvokerFor(desc Descriptor[T]) (Invoker[T], error)
	Update(arg A, v any)
	ToValue() V
}

type KeywordArg[T ref.References[T]] struct {
	ev  evaluate.Evaluator[T]
	arg string
	v   any
}

var _ Arg[ref.EmptyReferences] = &KeywordArg[ref.EmptyReferences]{}

func (ka *KeywordArg[T]) Expand(ctx context.Context, depth int) (*evaluate.Result[T], error) {
	v, found := ka.Get()
	if !found {
		return evaluate.EmptyResult[T](), nil
	}

	r, err := evaluateArg(ctx, ka.ev, ka.arg, v, depth)
	if err != nil {
		return nil, ka.WrapError(err)
	}

	return r, nil
}

func (ka *KeywordArg[_]) Get() (any, bool) {
	if ka == nil {
		return nil, false
	}

	return ka.v, true
}

func (ka *KeywordArg[_]) Annotation() evaluate.Annotation {
	return ArgAnnotation(ka.arg)
}

func (ka *KeywordArg[_]) WrapError(err error) error {
	return &KeywordArgError{Arg: ka.arg, Cause: err}
}

func MissingKeywordArg[T ref.References[T]]() *KeywordArg[T] {
	return nil
}

func NewKeywordArg[T ref.References[T]](ev evaluate.Evaluator[T], arg string, v any) *KeywordArg[T] {
	return &KeywordArg[T]{
		ev:  ev,
		arg: arg,
		v:   v,
	}
}

type KeywordArgs[T ref.References[T]] struct {
	ev   evaluate.Evaluator[T]
	args map[string]any
	into map[*KeywordArg[T]]*any
}

var _ Args[ref.EmptyReferences, *KeywordArg[ref.EmptyReferences], map[string]any] = &KeywordArgs[ref.EmptyReferences]{}

func (kas *KeywordArgs[T]) InvokerFor(desc Descriptor[T]) (Invoker[T], error) {
	return desc.KeywordInvoker(kas.ev, kas)
}

func (kas *KeywordArgs[T]) Update(arg *KeywordArg[T], v any) {
	if _, found := kas.into[arg]; found {
		kas.into[arg] = &v
	}
}

func (kas *KeywordArgs[T]) Use(k string) *KeywordArg[T] {
	v, found := kas.args[k]
	if !found {
		return MissingKeywordArg[T]()
	}

	arg := NewKeywordArg(kas.ev, k, v)
	kas.into[arg] = nil
	return arg
}

func (kas *KeywordArgs[T]) UseAll() map[string]*KeywordArg[T] {
	m := make(map[string]*KeywordArg[T], len(kas.args))
	for k := range kas.args {
		m[k] = kas.Use(k)
	}
	return m
}

func (kas *KeywordArgs[T]) Has(k string) bool {
	_, found := kas.args[k]
	return found
}

func (kas *KeywordArgs[T]) ToValue() map[string]any {
	out := make(map[string]any, len(kas.into))
	for arg, v := range kas.into {
		if v != nil {
			out[arg.arg] = *v
		} else {
			out[arg.arg] = kas.args[arg.arg]
		}
	}
	return out
}

func NewKeywordArgs[T ref.References[T]](ev evaluate.Evaluator[T], args map[string]any) *KeywordArgs[T] {
	return &KeywordArgs[T]{
		ev:   ev,
		args: args,
		into: make(map[*KeywordArg[T]]*any),
	}
}

type PositionalArg[T ref.References[T]] struct {
	ev  evaluate.Evaluator[T]
	arg int
	v   any
}

var _ Arg[ref.EmptyReferences] = &PositionalArg[ref.EmptyReferences]{}

func (pa *PositionalArg[T]) Expand(ctx context.Context, depth int) (*evaluate.Result[T], error) {
	v, found := pa.Get()
	if !found {
		return evaluate.EmptyResult[T](), nil
	}

	r, err := evaluateArg(ctx, pa.ev, pa.arg, v, depth)
	if err != nil {
		return nil, pa.WrapError(err)
	}

	return r, nil
}

func (pa *PositionalArg[_]) Get() (any, bool) {
	if pa == nil {
		return nil, false
	}

	return pa.v, true
}

func (pa *PositionalArg[_]) Annotation() evaluate.Annotation {
	return ArgAnnotation(pa.arg + 1)
}

func (pa *PositionalArg[_]) WrapError(err error) error {
	return &PositionalArgError{Arg: pa.arg + 1, Cause: err}
}

func MissingPositionalArg[T ref.References[T]]() *PositionalArg[T] {
	return nil
}

func NewPositionalArg[T ref.References[T]](ev evaluate.Evaluator[T], arg int, v any) *PositionalArg[T] {
	return &PositionalArg[T]{
		ev:  ev,
		arg: arg,
		v:   v,
	}
}

type PositionalArgs[T ref.References[T]] struct {
	ev   evaluate.Evaluator[T]
	args []any
	into map[*PositionalArg[T]]*any
}

var (
	_ Args[ref.EmptyReferences, *PositionalArg[ref.EmptyReferences], []any] = &PositionalArgs[ref.EmptyReferences]{}
)

func (pas *PositionalArgs[T]) InvokerFor(desc Descriptor[T]) (Invoker[T], error) {
	return desc.PositionalInvoker(pas.ev, pas)
}

func (pas *PositionalArgs[T]) Update(arg *PositionalArg[T], v any) {
	if _, found := pas.into[arg]; found {
		pas.into[arg] = &v
	}
}

func (pas *PositionalArgs[T]) Use(idx int) *PositionalArg[T] {
	if idx >= len(pas.args) {
		return MissingPositionalArg[T]()
	}

	arg := NewPositionalArg(pas.ev, idx, pas.args[idx])
	pas.into[arg] = nil
	return arg
}

func (pas *PositionalArgs[T]) UseAll() []*PositionalArg[T] {
	s := make([]*PositionalArg[T], len(pas.args))
	for idx := range pas.args {
		s[idx] = pas.Use(idx)
	}
	return s
}

func (pas *PositionalArgs[T]) Len() int {
	return len(pas.args)
}

func (pas *PositionalArgs[T]) ToValue() []any {
	out := make([]any, len(pas.args))
	for arg, v := range pas.into {
		if v != nil {
			out[arg.arg] = *v
		} else {
			out[arg.arg] = pas.args[arg.arg]
		}
	}
	return out
}

func NewPositionalArgs[T ref.References[T]](ev evaluate.Evaluator[T], args []any) *PositionalArgs[T] {
	return &PositionalArgs[T]{
		ev:   ev,
		args: args,
		into: make(map[*PositionalArg[T]]*any),
	}
}
