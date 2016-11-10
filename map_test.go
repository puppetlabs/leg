package godat

import (
	"errors"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapInsertionAndRetrieval(t *testing.T) {
	for _, m := range []Map{NewHashMap(), NewLinkedHashMap()} {
		assert.True(t, m.Empty())
		assert.False(t, m.Contains("a"))

		value, found := m.Get("a")
		assert.Nil(t, value)
		assert.False(t, found)

		found = m.Put("a", 1)
		assert.False(t, m.Empty())
		assert.False(t, found)

		value, found = m.Get("a")
		assert.Equal(t, 1, value)
		assert.True(t, found)
	}
}

func TestMapStoresNil(t *testing.T) {
	for _, m := range []Map{NewHashMap(), NewLinkedHashMap()} {
		m.Put("a", nil)

		assert.False(t, m.Empty())
		assert.True(t, m.Contains("a"))

		value, found := m.Get("a")
		assert.Nil(t, value)
		assert.True(t, found)

		var intoValue interface{} = 1
		found = m.GetInto("a", &intoValue)
		assert.Nil(t, intoValue)
		assert.True(t, found)
	}
}

func TestMapGetInto(t *testing.T) {
	for _, m := range []Map{NewHashMap(), NewLinkedHashMap()} {
		m.Put("a", 1)

		var value int
		found := m.GetInto("a", &value)
		assert.True(t, found)
		assert.Equal(t, 1, value)
	}
}

func TestMapUpdate(t *testing.T) {
	for _, m := range []Map{NewHashMap(), NewLinkedHashMap()} {
		m.Put("a", 1)
		m.Put("b", 2)

		var value int

		m.GetInto("a", &value)
		assert.Equal(t, 1, value)

		m.Put("a", 2)
		m.GetInto("a", &value)
		assert.Equal(t, 2, value)

		m.Remove("a")
		assert.False(t, m.Contains("a"))
		found := m.GetInto("a", &value)
		assert.Equal(t, 2, value)
		assert.False(t, found)

		m.Clear()
		assert.True(t, m.Empty())
	}
}

func TestMapIteration(t *testing.T) {
	for _, m := range []Map{NewHashMap(), NewLinkedHashMap(), NewHashMapWithCapacity(1), NewLinkedHashMapWithCapacity(1)} {
		m.Put("a", 1)
		m.Put("b", 2)
		m.Put("c", 3)

		assert.Equal(t, 3, m.Size())
		assert.Contains(t, m.Keys(), "a")
		assert.Contains(t, m.Keys(), "b")
		assert.Contains(t, m.Keys(), "c")
		assert.Contains(t, m.Values(), 1)
		assert.Contains(t, m.Values(), 2)
		assert.Contains(t, m.Values(), 3)

		err := m.ForEach(func(key, value interface{}) error {
			return errors.New("something went wrong!")
		})
		assert.EqualError(t, err, "something went wrong!")
	}
}

func TestMapIterationReflection(t *testing.T) {
	for _, m := range []Map{NewHashMap(), NewLinkedHashMap()} {
		m.Put("a", 1)
		m.Put("b", 2)
		m.Put("c", 3)

		var keysInto []string
		m.KeysInto(&keysInto)
		assert.Len(t, keysInto, 3)
		assert.Contains(t, keysInto, "a")
		assert.Contains(t, keysInto, "b")
		assert.Contains(t, keysInto, "c")

		var valuesInto []int
		m.ValuesInto(&valuesInto)
		assert.Len(t, valuesInto, 3)
		assert.Contains(t, valuesInto, 1)
		assert.Contains(t, valuesInto, 2)
		assert.Contains(t, valuesInto, 3)

		m.ForEachInto(func(key string, value int) error {
			assert.Contains(t, []string{"a", "b", "c"}, key)
			assert.Contains(t, []int{1, 2, 3}, value)

			return nil
		})

		err := m.ForEachInto(func(key string, value int) error {
			return errors.New("something else went wrong!")
		})
		assert.EqualError(t, err, "something else went wrong!")

		assert.Panics(t, func() {
			m.ForEachInto(func(key, value uint) error {
				return nil
			})
		})

		assert.Panics(t, func() {
			m.ForEachInto(func(key string, value int) (int, error) {
				return 0, nil
			})
		})
	}
}
