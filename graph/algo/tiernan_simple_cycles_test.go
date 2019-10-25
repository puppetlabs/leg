// Portions of this file are derived from JGraphT, a free Java graph-theory
// library.
//
// (C) Copyright 2013-2017, by Nikolay Ognyanov and Contributors.

package algo

import (
	"testing"

	"github.com/reflect/gographt"
	"github.com/stretchr/testify/assert"
)

func TestTiernanSimpleCycles(t *testing.T) {
	g := gographt.NewDirectedPseudographWithFeatures(gographt.DeterministicIteration)
	for i := 0; i < 7; i++ {
		g.AddVertex(i)
	}

	g.Connect(0, 0)
	assert.Equal(t, [][]gographt.Vertex{{0}}, TiernanSimpleCyclesOf(g).Cycles())

	g.Connect(1, 1)
	assert.Equal(t, [][]gographt.Vertex{{0}, {1}}, TiernanSimpleCyclesOf(g).Cycles())

	g.Connect(0, 1)
	g.Connect(1, 0)
	assert.Equal(t, [][]gographt.Vertex{{0, 1}, {0}, {1}}, TiernanSimpleCyclesOf(g).Cycles())

	g.Connect(1, 2)
	g.Connect(2, 3)
	g.Connect(3, 0)
	assert.Equal(t, [][]gographt.Vertex{{0, 1, 2, 3}, {0, 1}, {0}, {1}}, TiernanSimpleCyclesOf(g).Cycles())

	g.Connect(6, 6)
	assert.Equal(t, [][]gographt.Vertex{{0, 1, 2, 3}, {0, 1}, {0}, {1}, {6}}, TiernanSimpleCyclesOf(g).Cycles())

	conditions := []struct {
		ConnectedVertices, Cycles int
	}{
		{1, 1},
		{2, 3},
		{3, 8},
		{4, 24},
		{5, 89},
		{6, 415},
		{7, 2372},
		{8, 16072},
		{9, 125673},
	}

	for _, cond := range conditions {
		g = gographt.NewDirectedPseudograph()
		for i := 0; i < cond.ConnectedVertices; i++ {
			g.AddVertex(i)
		}
		for i := 0; i < cond.ConnectedVertices; i++ {
			for j := 0; j < cond.ConnectedVertices; j++ {
				g.Connect(i, j)
			}
		}

		assert.Len(t, TiernanSimpleCyclesOf(g).Cycles(), cond.Cycles)
	}
}

func TestTiernanSimpleCyclesInto(t *testing.T) {
	g := gographt.NewSimpleDirectedGraphWithFeatures(gographt.DeterministicIteration)
	g.AddVertex("a")
	g.AddVertex("b")
	g.AddVertex("c")

	g.Connect("a", "b")
	g.Connect("b", "c")
	g.Connect("c", "a")

	var cycles [][]string
	assert.NotPanics(t, func() {
		TiernanSimpleCyclesOf(g).CyclesInto(&cycles)
	})

	assert.Equal(t, [][]string{{"a", "b", "c"}}, cycles)
}
