// Portions of this file are derived from JGraphT, a free Java graph-theory
// library.
//
// (C) Copyright 2006-2016, by John V Sichi and Contributors.

package algo

import (
	"testing"

	"github.com/puppetlabs/horsehead/v2/graph"
	"github.com/stretchr/testify/assert"
)

func CreateShortestPathsTestGraph() graph.Graph {
	g := graph.NewSimpleWeightedGraph()

	g.AddVertex("A")
	g.AddVertex("B")
	g.AddVertex("C")
	g.AddVertex("D")
	g.AddVertex("E")

	g.AddEdgeWithWeight("A", "B", 1, 2.)
	g.AddEdgeWithWeight("A", "C", 2, 3.)
	g.AddEdgeWithWeight("B", "D", 3, 5.)
	g.AddEdgeWithWeight("C", "D", 4, 20.)
	g.AddEdgeWithWeight("D", "E", 5, 5.)
	g.AddEdgeWithWeight("A", "E", 6, 100.)

	return g
}

func TestBellmanFordShortestPaths(t *testing.T) {
	g := CreateShortestPathsTestGraph()
	bfsp := BellmanFordShortestPathsOf(g, "C")

	path, err := bfsp.EdgesTo("E")
	assert.NoError(t, err)
	assert.Equal(t, []graph.Edge{2, 1, 3, 5}, path)

	cost, err := bfsp.CostTo("E")
	assert.NoError(t, err)
	assert.Equal(t, 15., cost)
}
