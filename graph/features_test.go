package gographt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeterministicIterationFeature(t *testing.T) {
	g := NewSimpleGraphWithFeatures(DeterministicIteration)

	g.AddVertex("A")
	g.AddVertex("B")
	g.AddVertex("C")
	g.AddVertex("D")
	g.AddVertex("E")

	g.AddEdge("A", "B", 1)
	g.AddEdge("B", "C", 2)
	g.AddEdge("D", "E", 3)
	g.AddEdge("C", "D", 4)
	g.AddEdge("A", "E", 5)
	g.AddEdge("A", "D", 6)
	g.AddEdge("A", "C", 7)

	for i := 0; i < 10; i++ {
		assert.Equal(t, []Vertex{"A", "B", "C", "D", "E"}, g.Vertices().AsSlice())
		assert.Equal(t, []Edge{1, 2, 3, 4, 5, 6, 7}, g.Edges().AsSlice())

		edges, err := g.EdgesOf("A")
		assert.NoError(t, err)
		assert.Equal(t, []Edge{1, 5, 6, 7}, edges.AsSlice())
	}
}
