package gographt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUndirectedMultigraphEdgeAddition(t *testing.T) {
	g := NewUndirectedMultigraph()
	g.AddVertex("a")
	g.AddVertex("b")

	assert.NoError(t, g.AddEdge("a", "b", 1))

	assert.Equal(t, uint(2), g.Vertices().Count())
	assert.True(t, g.Vertices().Contains("a"))
	assert.True(t, g.Vertices().Contains("b"))

	edge, err := g.EdgeBetween("a", "b")
	assert.NoError(t, err)
	assert.Equal(t, 1, edge)

	assert.NoError(t, g.AddEdge("a", "b", 2))

	edges := g.EdgesBetween("a", "b")
	assert.Equal(t, uint(2), edges.Count())
	assert.True(t, edges.Contains(1))
	assert.True(t, edges.Contains(2))
}
