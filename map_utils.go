package godat

import (
	"reflect"
)

func mapGetInto(m Map, key, into interface{}) bool {
	value, found := m.Get(key)

	if found {
		t := reflect.ValueOf(into).Elem()
		rv := reflect.ValueOf(value)
		if !rv.IsValid() {
			rv = reflect.Zero(t.Type())
		}

		t.Set(rv)
	}

	return found
}

func mapKeys(m Map) []interface{} {
	keys := make([]interface{}, m.Size())

	i := 0
	m.ForEach(func(key, value interface{}) error {
		keys[i] = key
		i++

		return nil
	})

	return keys
}

func mapKeysInto(m Map, into interface{}) {
	p := reflect.ValueOf(into).Elem()
	slice := p

	m.ForEach(func(key, value interface{}) error {
		slice = reflect.Append(slice, reflect.ValueOf(key))
		return nil
	})

	p.Set(slice)
}

func mapValues(m Map) []interface{} {
	values := make([]interface{}, m.Size())

	i := 0
	m.ForEach(func(key, value interface{}) error {
		values[i] = value
		i++

		return nil
	})

	return values
}

func mapValuesInto(m Map, into interface{}) {
	p := reflect.ValueOf(into).Elem()
	slice := p

	m.ForEach(func(key, value interface{}) error {
		slice = reflect.Append(slice, reflect.ValueOf(value))
		return nil
	})

	p.Set(slice)
}

func mapForEachInto(m Map, fn interface{}) error {
	fnr := reflect.ValueOf(fn)

	return m.ForEach(func(key, value interface{}) error {
		if fnr.Type().NumOut() != 1 {
			panic(ErrInvalidFuncSignature)
		}

		r := fnr.Call([]reflect.Value{
			reflect.ValueOf(key),
			reflect.ValueOf(value),
		})

		err := r[0]
		if err.IsNil() {
			return nil
		} else {
			return err.Interface().(error)
		}
	})
}
