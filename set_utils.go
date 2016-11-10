package godat

import (
	"reflect"
)

func setValuesInto(s Set, into interface{}) {
	p := reflect.ValueOf(into).Elem()
	slice := p

	s.ForEach(func(element interface{}) error {
		slice = reflect.Append(slice, reflect.ValueOf(element))
		return nil
	})

	p.Set(slice)
}

func setForEachInto(s Set, fn interface{}) error {
	fnr := reflect.ValueOf(fn)

	return s.ForEach(func(element interface{}) error {
		if fnr.Type().NumOut() != 1 {
			panic(ErrInvalidFuncSignature)
		}

		r := fnr.Call([]reflect.Value{reflect.ValueOf(element)})

		err := r[0]
		if err.IsNil() {
			return nil
		} else {
			return err.Interface().(error)
		}
	})
}
