package godat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	pq := NewPriorityQueue()
	assert.True(t, pq.Empty())

	pq.Add("a", 10)
	pq.Add("c", 30)
	pq.Add("b", 20)

	assert.Equal(t, 3, pq.Size())

	value, found := pq.Poll()
	assert.Equal(t, "c", value)
	assert.True(t, found)

	value, _ = pq.Poll()
	assert.Equal(t, "b", value)

	pq.Add("x", 5)

	value, _ = pq.Poll()
	assert.Equal(t, "a", value)

	pq.Clear()
	assert.True(t, pq.Empty())
}

func TestPriorityQueueReflection(t *testing.T) {
	pq := NewPriorityQueue()

	pq.Add("a", 10)
	pq.Add("b", 20)
	pq.Add(nil, 15)

	var value string
	found := pq.PollInto(&value)
	assert.Equal(t, "b", value)
	assert.True(t, found)

	var intoNilValue interface{} = 1
	found = pq.PollInto(&intoNilValue)
	assert.Equal(t, nil, intoNilValue)
	assert.True(t, found)

	pq.PollInto(&value)
	assert.Equal(t, "a", value)

	found = pq.PollInto(&value)
	assert.Equal(t, "a", value)
	assert.False(t, found)
}
