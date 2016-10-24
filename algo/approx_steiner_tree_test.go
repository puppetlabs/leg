package algo

import (
	"github.com/reflect/gographt"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApproximateSteinerTree(t *testing.T) {
	g := gographt.NewUndirectedPseudograph()

	g.AddVertex("A")
	g.AddVertex("B")
	g.AddVertex("C")
	g.AddVertex("D")
	g.AddVertex("E")
	g.AddVertex("F")
	g.AddVertex("G")

	g.AddEdge("A", "B", 1)
	g.AddEdge("B", "C", 2)
	g.AddEdge("C", "E", 3)
	g.AddEdge("C", "F", 4)
	g.AddEdge("D", "A", 5)

	reduced, err := ApproximateSteinerTreeOf(g, []gographt.Vertex{"A", "D", "F"})
	assert.NoError(t, err)
	assert.Equal(t, uint(5), reduced.Vertices().Count())

	edges := reduced.Edges().AsSlice()
	assert.Contains(t, edges, 1)
	assert.Contains(t, edges, 2)
	assert.Contains(t, edges, 4)
	assert.Contains(t, edges, 5)
}
