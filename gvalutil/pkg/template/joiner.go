package template

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/PaesslerAG/gval"
)

type Joiner interface {
	Join(ctx context.Context, a, b gval.Evaluable, parameter interface{}) (interface{}, error)
}

type JoinerFunc func(ctx context.Context, a, b gval.Evaluable, parameter interface{}) (interface{}, error)

var _ Joiner = JoinerFunc(nil)

func (jf JoinerFunc) Join(ctx context.Context, a, b gval.Evaluable, parameter interface{}) (interface{}, error) {
	return jf(ctx, a, b, parameter)
}

type StringFormatter interface {
	FormatString(v interface{}) (string, error)
}

type StringFormatterFunc func(v interface{}) (string, error)

var _ StringFormatter = StringFormatterFunc(nil)

func (sff StringFormatterFunc) FormatString(v interface{}) (string, error) {
	return sff(v)
}

type defaultStringFormatter struct{}

func (dsf defaultStringFormatter) FormatString(v interface{}) (string, error) {
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

func (sj *StringJoiner) Join(ctx context.Context, a, b gval.Evaluable, parameter interface{}) (interface{}, error) {
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

	fa, err := sj.formatter.FormatString(ea)
	if err != nil {
		return nil, err
	}

	fb, err := sj.formatter.FormatString(eb)
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
