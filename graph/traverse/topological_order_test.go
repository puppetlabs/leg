// Portions of this file are derived from JGraphT, a free Java graph-theory
// library.
//
// (C) Copyright 2005-2017, by John V Sichi and Contributors.

package traverse

import (
	"math/rand"
	"testing"

	"github.com/puppetlabs/horsehead/v2/graph"
	"github.com/stretchr/testify/assert"
)

func TestTopologicalOrderTraverser(t *testing.T) {
	// This is really the test case from the JGraphT source code.
	v := []string{
		"preheat oven",
		"sift dry ingredients",
		"stir wet ingredients",
		"mix wet and dry ingredients",
		"spoon onto pan",
		"bake",
		"cool",
		"frost",
		"eat",
	}

	g := graph.NewSimpleDirectedGraph()
	for _, i := range rand.Perm(len(v)) {
		g.AddVertex(v[i])
	}

	g.Connect(v[0], v[1])
	g.Connect(v[1], v[2])
	g.Connect(v[0], v[2])
	g.Connect(v[1], v[3])
	g.Connect(v[2], v[3])
	g.Connect(v[3], v[4])
	g.Connect(v[4], v[5])
	g.Connect(v[5], v[6])
	g.Connect(v[6], v[7])
	g.Connect(v[7], v[8])
	g.Connect(v[6], v[8])

	i := 0
	err := NewTopologicalOrderTraverser(g).ForEach(func(vertex graph.Vertex) error {
		assert.Equal(t, v[i], vertex)
		i++

		return nil
	})
	assert.NoError(t, err)
	assert.Len(t, v, i)
}

func TestTopologicalOrderTraverserWithDisconnectedRoots(t *testing.T) {
	g := graph.NewSimpleDirectedGraphWithFeatures(graph.DeterministicIteration)

	g.AddVertex("a")
	g.AddVertex("b")
	g.Connect("a", "b")

	g.AddVertex(1)
	g.AddVertex(2)
	g.Connect(1, 2)

	vertices, err := NewTopologicalOrderTraverser(g).Vertices()
	assert.NoError(t, err)
	assert.Equal(t, []graph.Vertex{"a", 1, "b", 2}, vertices)
}

func TestTopologicalOrderTraverserForEachInto(t *testing.T) {
	g := graph.NewSimpleDirectedGraph()

	v := []string{"a", "b", "c"}
	for _, i := range rand.Perm(len(v)) {
		g.AddVertex(v[i])
	}

	g.Connect(v[0], v[1])
	g.Connect(v[1], v[2])

	i := 0
	err := NewTopologicalOrderTraverser(g).ForEachInto(func(vertex string) error {
		assert.Equal(t, v[i], vertex)
		i++

		return nil
	})
	assert.NoError(t, err)
	assert.Len(t, v, i)
}
