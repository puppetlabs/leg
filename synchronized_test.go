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
