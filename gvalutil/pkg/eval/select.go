package eval

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/PaesslerAG/gval"
)

func IntValue(val any) (int, error) {
	switch vt := val.(type) {
	case string:
		vi, err := strconv.ParseInt(vt, 10, 32)
		if err != nil {
			return 0, &UnexpectedStringIndexError{RawIndex: vt, Cause: err}
		}

		return int(vi), nil
	case int, int8, int16, int32, int64:
		return int(reflect.ValueOf(vt).Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return int(reflect.ValueOf(vt).Uint()), nil
	case float32, float64:
		return int(reflect.ValueOf(vt).Float()), nil
	default:
		return 0, &UnexpectedIndexTypeError{RawIndex: vt}
	}

}

func SelectIndex(parameter []any, idx any) (offset int, value any, err error) {
	i, err := IntValue(idx)
	if err != nil {
		return 0, nil, err
	}

	if i < 0 && len(parameter)+i >= 0 {
		i = len(parameter) + i
	} else if i < 0 || i >= len(parameter) {
		return 0, nil, &IndexOutOfBoundsError{Index: i}
	}

	return i, parameter[i], nil
}

func StringValue(val any) (string, error) {
	switch vt := val.(type) {
	case fmt.Stringer:
		return vt.String(), nil
	case string:
		return vt, nil
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(reflect.ValueOf(vt).Int(), 10), nil
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(reflect.ValueOf(vt).Uint(), 10), nil
	case float32, float64:
		return strconv.FormatFloat(reflect.ValueOf(vt).Float(), 'f', -1, 64), nil
	default:
		return "", &UnexpectedKeyTypeError{RawKey: vt}
	}
}

func SelectKey(parameter map[string]any, key any) (name string, value any, err error) {
	k, err := StringValue(key)
	if err != nil {
		return "", nil, err
	}

	r, ok := parameter[k]
	if !ok {
		return "", nil, &UnknownKeyError{Key: k}
	}

	return k, r, nil
}

func SelectField(parameter, field any) (name string, value any, err error) {
	f, err := StringValue(field)
	if err != nil {
		return "", nil, err
	}

	pv := reflect.ValueOf(parameter)
	for {
		switch pv.Kind() {
		case reflect.Ptr:
			pv = pv.Elem()
		case reflect.Struct:
			fv := pv.FieldByName(f)
			if fv.IsValid() && fv.CanInterface() {
				return f, fv.Interface(), nil
			}

			return "", nil, &UnknownFieldError{Field: f}
		default:
			return "", nil, &UnsupportedValueTypeError{Value: parameter, Field: f}
		}
	}
}

func Select(ctx context.Context, parameter, elem any) (out any, err error) {
	switch pt := parameter.(type) {
	case Indexable:
		out, err = pt.Index(ctx, elem)
	case []any:
		_, out, err = SelectIndex(pt, elem)
	case map[string]any:
		_, out, err = SelectKey(pt, elem)
	default:
		_, out, err = SelectField(pt, elem)
	}
	return
}

func Selector(path gval.Evaluables) gval.Evaluable {
	return func(ctx context.Context, parameter any) (out any, err error) {
		out = parameter
		for _, eval := range path {
			elem, err := eval(ctx, parameter)
			if err != nil {
				return nil, err
			}

			out, err = Select(ctx, out, elem)
			if err != nil {
				return nil, err
			}
		}
		return
	}
}
