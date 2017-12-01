package gographt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOppositeVertexOf(t *testing.T) {
	g := NewSimpleGraph()
	g.AddVertex("a")
	g.AddVertex("b")

	assert.NoError(t, g.Connect("a", "b"))
	edges, err := g.EdgesOf("a")
	assert.NoError(t, err)
	assert.Equal(t, edges.Count(), uint(1))

	v, err := OppositeVertexOf(g, edges.AsSlice()[0], "a")
	assert.NoError(t, err)
	assert.Equal(t, "b", v)

	v, err = OppositeVertexOf(g, edges.AsSlice()[0], v)
	assert.NoError(t, err)
	assert.Equal(t, "a", v)
}
