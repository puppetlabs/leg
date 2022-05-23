package template

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/PaesslerAG/gval"
)

type Joiner interface {
	Join(ctx context.Context, a, b gval.Evaluable, parameter any) (any, error)
}

type JoinerFunc func(ctx context.Context, a, b gval.Evaluable, parameter any) (any, error)

var _ Joiner = JoinerFunc(nil)

func (jf JoinerFunc) Join(ctx context.Context, a, b gval.Evaluable, parameter any) (any, error) {
	return jf(ctx, a, b, parameter)
}

type StringFormatter interface {
	FormatString(ctx context.Context, v any) (string, error)
}

type StringFormatterFunc func(ctx context.Context, v any) (string, error)

var _ StringFormatter = StringFormatterFunc(nil)

func (sff StringFormatterFunc) FormatString(ctx context.Context, v any) (string, error) {
	return sff(ctx, v)
}

type defaultStringFormatter struct{}

func (dsf defaultStringFormatter) FormatString(ctx context.Context, v any) (string, error) {
	switch vt := v.(type) {
	case fmt.Stringer:
		return vt.String(), nil
	case string:
		return vt, nil
	case []byte:
		return string(vt), nil
	default:
		b, err := json.MarshalIndent(v, "", "\t")
		if err != nil {
			return "", err
		}

		return string(b), nil
	}
}

var DefaultStringFormatter StringFormatter = &defaultStringFormatter{}

type StringJoiner struct {
	formatter             StringFormatter
	eliminateEmptyStrings bool
}

var _ Joiner = &StringJoiner{}

func (sj *StringJoiner) Join(ctx context.Context, a, b gval.Evaluable, parameter any) (any, error) {
	ea, err := a(ctx, parameter)
	if err != nil {
		return nil, err
	}

	eb, err := b(ctx, parameter)
	if err != nil {
		return nil, err
	}

	if sj.eliminateEmptyStrings {
		switch {
		case ea == "":
			return eb, nil
		case eb == "":
			return ea, nil
		}
	}

	fa, err := sj.formatter.FormatString(ctx, ea)
	if err != nil {
		return nil, err
	}

	fb, err := sj.formatter.FormatString(ctx, eb)
	if err != nil {
		return nil, err
	}

	return fa + fb, nil
}

type StringJoinerOptions struct {
	Formatter             StringFormatter
	EliminateEmptyStrings bool
}

type StringJoinerOption interface {
	ApplyToStringJoinerOptions(target *StringJoinerOptions)
}

func (o *StringJoinerOptions) ApplyOptions(opts []StringJoinerOption) {
	for _, opt := range opts {
		opt.ApplyToStringJoinerOptions(o)
	}
}

func NewStringJoiner(opts ...StringJoinerOption) *StringJoiner {
	o := &StringJoinerOptions{
		Formatter: DefaultStringFormatter,
	}
	o.ApplyOptions(opts)

	return &StringJoiner{
		formatter:             o.Formatter,
		eliminateEmptyStrings: o.EliminateEmptyStrings,
	}
}

var DefaultStringJoiner = NewStringJoiner()
