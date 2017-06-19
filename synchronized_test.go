package godat

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSynchronizedContainersDoNotStack(t *testing.T) {
	sm1 := NewSynchronizedMap(NewHashMap())
	sm2 := NewSynchronizedMap(sm1)

	assert.Equal(t, reflect.ValueOf(sm1).Pointer(), reflect.ValueOf(sm2).Pointer())

	ss1 := NewSynchronizedSet(NewHashSet())
	ss2 := NewSynchronizedSet(ss1)

	assert.Equal(t, reflect.ValueOf(ss1).Pointer(), reflect.ValueOf(ss2).Pointer())
}

func TestSynchronizedMapAtomicOperations(t *testing.T) {
	sm1 := NewSynchronizedMap(NewHashMap())
	sm1.Put("a", 1)

	assert.True(t, sm1.CompareAndPut("a", 2, 1))

	v, found := sm1.Get("a")
	assert.True(t, found)
	assert.Equal(t, 2, v)

	assert.False(t, sm1.CompareAndPut("a", 3, 1))

	v, found = sm1.Get("a")
	assert.True(t, found)
	assert.Equal(t, 2, v)

	assert.False(t, sm1.CompareAndRemove("a", 1))

	v, found = sm1.Get("a")
	assert.True(t, found)
	assert.Equal(t, 2, v)

	assert.True(t, sm1.CompareAndRemove("a", 2))

	_, found = sm1.Get("a")
	assert.False(t, found)
}
