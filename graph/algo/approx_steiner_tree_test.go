package algo

import (
	"github.com/reflect/gographt"

	"math/rand"

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
	assert.Equal(t, uint(5), reduced.AsGraph().Vertices().Count())

	edges := reduced.Edges().AsSlice()
	assert.Contains(t, edges, 1)
	assert.Contains(t, edges, 2)
	assert.Contains(t, edges, 4)
	assert.Contains(t, edges, 5)
}

func TestApproximateSteinerTreeOfDeterministicGraph(t *testing.T) {
	g := gographt.NewUndirectedPseudographWithFeatures(gographt.DeterministicIteration)

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
	g.AddEdge("C", "G", 6)

	required := []gographt.Vertex{"A", "D", "F"}

	for i := 0; i < 10; i++ {
		shuffled := make([]gographt.Vertex, len(required))
		for in, out := range rand.Perm(len(required)) {
			shuffled[out] = required[in]
		}

		reduced, err := ApproximateSteinerTreeOf(g, shuffled)
		assert.NoError(t, err)
		assert.Equal(t, []gographt.Vertex{"A", "D", "F", "B", "C"}, reduced.AsGraph().Vertices().AsSlice())
		assert.Equal(t, []gographt.Edge{5, 1, 2, 4}, reduced.Edges().AsSlice())
	}
}
