package eval

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/PaesslerAG/gval"
)

func IntValue(val interface{}) (int, error) {
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

func SelectIndex(parameter []interface{}, idx interface{}) (offset int, value interface{}, err error) {
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

func StringValue(val interface{}) (string, error) {
	switch vt := val.(type) {
	case fmt.Stringer:
		return vt.String(), nil
	case string:
		return vt, nil
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(reflect.ValueOf(vt).Int(), 64), nil
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(reflect.ValueOf(vt).Uint(), 64), nil
	case float32, float64:
		return strconv.FormatFloat(reflect.ValueOf(vt).Float(), 'f', -1, 64), nil
	default:
		return "", &UnexpectedKeyTypeError{RawKey: vt}
	}
}

func SelectKey(parameter map[string]interface{}, key interface{}) (name string, value interface{}, err error) {
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

func SelectField(parameter, field interface{}) (name string, value interface{}, err error) {
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

func Select(ctx context.Context, parameter, elem interface{}) (out interface{}, err error) {
	switch pt := parameter.(type) {
	case Indexable:
		out, err = pt.Index(ctx, elem)
	case []interface{}:
		_, out, err = SelectIndex(pt, elem)
	case map[string]interface{}:
		_, out, err = SelectKey(pt, elem)
	default:
		_, out, err = SelectField(pt, elem)
	}
	return
}

func Selector(path gval.Evaluables) gval.Evaluable {
	return func(ctx context.Context, parameter interface{}) (out interface{}, err error) {
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
