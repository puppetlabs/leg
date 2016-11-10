package godat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLinkedHashMapIterationOrder(t *testing.T) {
	m := NewLinkedHashMap()

	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)

	keys := []interface{}{"a", "b", "c"}
	values := []interface{}{1, 2, 3}

	for ti := 0; ti < 10; ti++ {
		// Make sure the order stays the same!
		m.Put("b", 2)

		assert.Equal(t, keys, m.Keys())
		assert.Equal(t, values, m.Values())

		i := 0
		m.ForEach(func(key, value interface{}) error {
			assert.Equal(t, keys[i], key)
			assert.Equal(t, values[i], value)
			i++

			return nil
		})
	}
}
