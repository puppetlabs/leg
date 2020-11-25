// Portions of this file are derived from JGraphT, a free Java graph-theory
// library.
//
// (C) Copyright 2010-2016, by Tom Conerly and Contributors.

package algo

import (
	"testing"

	"github.com/puppetlabs/leg/graph"
	"github.com/stretchr/testify/assert"
)

func TestPrimWithConnectedWeightedGraph(t *testing.T) {
	g := graph.NewSimpleWeightedGraph()

	for _, vertex := range []string{"A", "B", "C", "D", "E"} {
		g.AddVertex(vertex)
	}

	// The resulting edges should be:
	//   1, 2, 3, 5
	//
	// The resulting total weight should be:
	//   2 + 3 + 5 + 5 = 15

	assert.NoError(t, g.AddEdgeWithWeight("A", "B", 1, 2.))
	assert.NoError(t, g.AddEdgeWithWeight("A", "C", 2, 3.))
	assert.NoError(t, g.AddEdgeWithWeight("B", "D", 3, 5.))
	assert.NoError(t, g.AddEdgeWithWeight("C", "D", 4, 20.))
	assert.NoError(t, g.AddEdgeWithWeight("D", "E", 5, 5.))
	assert.NoError(t, g.AddEdgeWithWeight("A", "E", 6, 100.))

	mst := PrimMinimumSpanningTreeOf(g)
	assert.Equal(t, uint(4), mst.Edges().Count())

	edges := mst.Edges().AsSlice()
	assert.Contains(t, edges, 1)
	assert.Contains(t, edges, 2)
	assert.Contains(t, edges, 3)
	assert.Contains(t, edges, 5)
	assert.InDelta(t, 15., mst.TotalWeight, 1e-6)
}

func TestPrimWithEqualWeights(t *testing.T) {
	g := graph.NewSimpleGraph()

	for _, vertex := range []string{"A", "B", "C", "D", "E"} {
		g.AddVertex(vertex)
	}

	assert.NoError(t, g.Connect("A", "B"))
	assert.NoError(t, g.Connect("A", "C"))
	assert.NoError(t, g.Connect("B", "D"))
	assert.NoError(t, g.Connect("C", "D"))
	assert.NoError(t, g.Connect("D", "E"))
	assert.NoError(t, g.Connect("A", "E"))

	mst := PrimMinimumSpanningTreeOf(g)
	assert.Equal(t, uint(4), mst.Edges().Count())

	edges := mst.Edges().AsSlice()
	assert.Len(t, edges, 4)
	assert.InDelta(t, graph.DefaultEdgeWeight*4, mst.TotalWeight, 1e-6)
}
