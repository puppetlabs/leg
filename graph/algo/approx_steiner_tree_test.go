package algo

import (
	"math/rand"
	"testing"

	"github.com/puppetlabs/horsehead/v2/graph"
	"github.com/stretchr/testify/assert"
)

func TestApproximateSteinerTree(t *testing.T) {
	g := graph.NewUndirectedPseudograph()

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

	reduced, err := ApproximateSteinerTreeOf(g, []graph.Vertex{"A", "D", "F"})
	assert.NoError(t, err)
	assert.Equal(t, uint(5), reduced.AsGraph().Vertices().Count())

	edges := reduced.Edges().AsSlice()
	assert.Contains(t, edges, 1)
	assert.Contains(t, edges, 2)
	assert.Contains(t, edges, 4)
	assert.Contains(t, edges, 5)
}

func TestApproximateSteinerTreeOfDeterministicGraph(t *testing.T) {
	g := graph.NewUndirectedPseudographWithFeatures(graph.DeterministicIteration)

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

	required := []graph.Vertex{"A", "D", "F"}

	for i := 0; i < 10; i++ {
		shuffled := make([]graph.Vertex, len(required))
		for in, out := range rand.Perm(len(required)) {
			shuffled[out] = required[in]
		}

		reduced, err := ApproximateSteinerTreeOf(g, shuffled)
		assert.NoError(t, err)
		assert.Equal(t, []graph.Vertex{"A", "D", "F", "B", "C"}, reduced.AsGraph().Vertices().AsSlice())
		assert.Equal(t, []graph.Edge{5, 1, 2, 4}, reduced.Edges().AsSlice())
	}
}
