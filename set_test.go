package godat

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetInsertionAndRetrieval(t *testing.T) {
	for _, s := range []Set{NewHashSet(), NewLinkedHashSet(), NewSynchronizedSet(NewHashSet())} {
		assert.True(t, s.Empty())
		assert.False(t, s.Contains("a", "b"))

		s.Add("a", "b", "c")
		assert.Equal(t, 3, s.Size())
		assert.True(t, s.Contains("a"))
		assert.True(t, s.Contains("a", "b", "c"))
		assert.False(t, s.Contains("d"))

		s.Add("c")
		assert.Equal(t, 3, s.Size())
	}
}

func TestSetRemoval(t *testing.T) {
	for _, s := range []Set{NewHashSet(), NewLinkedHashSet(), NewSynchronizedSet(NewHashSet())} {
		s.Add("a", "b", "c")
		s.Remove("a", "b")

		assert.Equal(t, 1, s.Size())
		assert.True(t, s.Contains("c"))
		assert.False(t, s.Contains("a"))

		s.Clear()
		assert.True(t, s.Empty())
	}
}

func TestSetIteration(t *testing.T) {
	for _, s := range []Set{NewHashSet(), NewLinkedHashSet(), NewHashSetWithCapacity(1), NewLinkedHashSetWithCapacity(1), NewSynchronizedSet(NewHashSet())} {
		s.Add("a", "b", "c")

		assert.Equal(t, 3, s.Size())
		assert.Contains(t, s.Values(), "a")
		assert.Contains(t, s.Values(), "b")
		assert.Contains(t, s.Values(), "c")

		err := s.ForEach(func(element interface{}) error {
			return errors.New("something went wrong!")
		})
		assert.EqualError(t, err, "something went wrong!")
	}
}

func TestLinkedHashSetIteration(t *testing.T) {
	values := []interface{}{"a", "b", "c"}

	s := NewLinkedHashSet()
	s.Add(values...)

	for ti := 0; ti < 10; ti++ {
		assert.Equal(t, values, s.Values())

		i := 0
		s.ForEach(func(element interface{}) error {
			assert.Equal(t, values[i], element)
			i++

			return nil
		})
	}
}

func TestSetIterationReflection(t *testing.T) {
	for _, s := range []Set{NewHashSet(), NewLinkedHashSet(), NewSynchronizedSet(NewHashSet())} {
		s.Add("a", "b", "c")

		var valuesInto []string
		s.ValuesInto(&valuesInto)
		assert.Len(t, valuesInto, 3)
		assert.Contains(t, valuesInto, "a")
		assert.Contains(t, valuesInto, "b")
		assert.Contains(t, valuesInto, "c")

		s.ForEachInto(func(element string) error {
			assert.Contains(t, []string{"a", "b", "c"}, element)

			return nil
		})

		err := s.ForEachInto(func(element string) error {
			return errors.New("something else went wrong!")
		})
		assert.EqualError(t, err, "something else went wrong!")

		assert.Panics(t, func() {
			s.ForEachInto(func(element uint) error {
				return nil
			})
		})

		assert.Panics(t, func() {
			s.ForEachInto(func(element string) (int, error) {
				return 0, nil
			})
		})
	}
}

func TestSetOperations(t *testing.T) {
	for _, s := range []Set{NewHashSet(), NewLinkedHashSet(), NewSynchronizedSet(NewHashSet())} {
		s.Add("a", "b", "c")

		u := NewHashSet()
		u.Add("c", "d", "e")

		s.AddAll(u)
		assert.Equal(t, 5, s.Size())
		assert.Contains(t, s.Values(), "a")
		assert.Contains(t, s.Values(), "b")
		assert.Contains(t, s.Values(), "c")
		assert.Contains(t, s.Values(), "d")
		assert.Contains(t, s.Values(), "e")

		rn := NewHashSet()
		rn.Add("a", "c", "e", "f")

		s.RetainAll(rn)
		assert.Equal(t, 3, s.Size())
		assert.Contains(t, s.Values(), "a")
		assert.Contains(t, s.Values(), "c")
		assert.Contains(t, s.Values(), "e")

		rm := NewHashSet()
		rm.Add("d", "e", "f")

		s.RemoveAll(rm)
		assert.Equal(t, 2, s.Size())
		assert.Contains(t, s.Values(), "a")
		assert.Contains(t, s.Values(), "c")
	}
}
