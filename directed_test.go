package gographt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	SimpleDirectedGraphConstructors = []func() MutableDirectedGraph{
		func() MutableDirectedGraph { return NewSimpleDirectedGraph() },
		func() MutableDirectedGraph { return NewSimpleDirectedGraphWithFeatures(DeterministicIteration) },
		func() MutableDirectedGraph { return NewSimpleDirectedWeightedGraph() },
		func() MutableDirectedGraph { return NewSimpleDirectedWeightedGraphWithFeatures(DeterministicIteration) },
	}

	DirectedMultigraphConstructors = []func() MutableDirectedGraph{
		func() MutableDirectedGraph { return NewDirectedMultigraph() },
		func() MutableDirectedGraph { return NewDirectedMultigraphWithFeatures(DeterministicIteration) },
		func() MutableDirectedGraph { return NewDirectedWeightedMultigraph() },
		func() MutableDirectedGraph { return NewDirectedWeightedMultigraphWithFeatures(DeterministicIteration) },
	}

	DirectedPseudographConstructors = []func() MutableDirectedGraph{
		func() MutableDirectedGraph { return NewDirectedPseudograph() },
		func() MutableDirectedGraph { return NewDirectedPseudographWithFeatures(DeterministicIteration) },
		func() MutableDirectedGraph { return NewDirectedWeightedPseudograph() },
		func() MutableDirectedGraph { return NewDirectedWeightedPseudographWithFeatures(DeterministicIteration) },
	}

	DirectedConstructors = append(
		append(
			append([]func() MutableDirectedGraph{}, SimpleDirectedGraphConstructors...),
			DirectedMultigraphConstructors...),
		DirectedPseudographConstructors...)
)

func TestDirectedMultigraphEdgeAddition(t *testing.T) {
	for _, constructor := range DirectedMultigraphConstructors {
		g := constructor()
		g.AddVertex("a")
		g.AddVertex("b")

		assert.NoError(t, g.AddEdge("a", "b", 1))

		assert.Equal(t, uint(2), g.Vertices().Count())
		assert.True(t, g.Vertices().Contains("a"))
		assert.True(t, g.Vertices().Contains("b"))

		edge, err := g.EdgeBetween("a", "b")
		assert.NoError(t, err)
		assert.Equal(t, 1, edge)

		edge, err = g.EdgeBetween("b", "a")
		assert.Equal(t, &NotConnectedError{"b", "a"}, err)
		assert.Nil(t, edge)

		assert.NoError(t, g.AddEdge("a", "b", 2))

		edges := g.EdgesBetween("a", "b")
		assert.Equal(t, uint(2), edges.Count())
		assert.True(t, edges.Contains(1))
		assert.True(t, edges.Contains(2))

		in, err := g.InDegreeOf("a")
		assert.NoError(t, err)
		assert.Equal(t, uint(0), in)

		out, err := g.OutDegreeOf("a")
		assert.NoError(t, err)
		assert.Equal(t, uint(2), out)
	}
}

func TestDirectedConnections(t *testing.T) {
	for _, constructor := range DirectedConstructors {
		g := constructor()
		g.AddVertex("a")
		g.AddVertex("b")
		g.AddVertex("c")

		assert.NoError(t, g.AddEdge("a", "b", 1))
		assert.NoError(t, g.AddEdge("b", "a", 2))
		assert.NoError(t, g.AddEdge("b", "c", 3))

		counts := map[string]*struct{ In, Out uint }{
			"a": {In: 1, Out: 1},
			"b": {In: 1, Out: 2},
			"c": {In: 1, Out: 0},
		}

		for vertex, expected := range counts {
			n, err := g.InDegreeOf(vertex)
			assert.NoError(t, err)
			assert.Equal(t, expected.In, n, "for vertex %s", vertex)

			n, err = g.OutDegreeOf(vertex)
			assert.NoError(t, err)
			assert.Equal(t, expected.Out, n, "for vertex %s", vertex)
		}

		assert.True(t, g.RemoveEdge(1))

		// We can't remove the same edge again, obviously.
		assert.False(t, g.RemoveEdge(1))

		counts["a"].Out--
		counts["b"].In--

		for vertex, expected := range counts {
			n, err := g.InDegreeOf(vertex)
			assert.NoError(t, err)
			assert.Equal(t, expected.In, n, "for vertex %s", vertex)

			n, err = g.OutDegreeOf(vertex)
			assert.NoError(t, err)
			assert.Equal(t, expected.Out, n, "for vertex %s", vertex)
		}

		edge, err := g.RemoveEdgeBetween("b", "a")
		assert.NoError(t, err)
		assert.Equal(t, 2, edge)

		counts["a"].In--
		counts["b"].Out--

		for vertex, expected := range counts {
			n, err := g.InDegreeOf(vertex)
			assert.NoError(t, err)
			assert.Equal(t, expected.In, n, "for vertex %s", vertex)

			n, err = g.OutDegreeOf(vertex)
			assert.NoError(t, err)
			assert.Equal(t, expected.Out, n, "for vertex %s", vertex)
		}

		assert.True(t, g.RemoveVertex("c"))

		// We can't remove the same vertex either.
		assert.False(t, g.RemoveVertex("c"))

		delete(counts, "c")
		counts["b"].Out--

		for vertex, expected := range counts {
			n, err := g.InDegreeOf(vertex)
			assert.NoError(t, err)
			assert.Equal(t, expected.In, n, "for vertex %s", vertex)

			n, err = g.OutDegreeOf(vertex)
			assert.NoError(t, err)
			assert.Equal(t, expected.Out, n, "for vertex %s", vertex)
		}

		_, err = g.InDegreeOf("c")
		assert.Equal(t, &VertexNotFoundError{Vertex: "c"}, err)

		_, err = g.OutDegreeOf("c")
		assert.Equal(t, &VertexNotFoundError{Vertex: "c"}, err)
	}
}

func TestDirectedVertexNotFound(t *testing.T) {
	for _, constructor := range DirectedConstructors {
		g := constructor()
		g.AddVertex("a")
		g.AddVertex("b")
		g.Connect("a", "b")

		edge, err := g.EdgeBetween("a", "c")
		assert.Nil(t, edge)
		assert.Equal(t, &NotConnectedError{Source: "a", Target: "c"}, err)

		edge, err = g.EdgesOf("c")
		assert.Nil(t, edge)
		assert.Equal(t, &VertexNotFoundError{Vertex: "c"}, err)
	}
}

func TestDirectedEdgeNotFound(t *testing.T) {
	for _, constructor := range DirectedConstructors {
		g := constructor()
		g.AddVertex("a")
		g.AddVertex("b")
		g.AddVertex("c")
		g.Connect("b", "c")

		edge, err := g.EdgeBetween("a", "b")
		assert.Nil(t, edge)
		assert.Equal(t, &NotConnectedError{Source: "a", Target: "b"}, err)

		edge, err = g.EdgeBetween("c", "b")
		assert.Nil(t, edge)
		assert.Equal(t, &NotConnectedError{Source: "c", Target: "b"}, err)

		v, err := g.SourceVertexOf(NewEdge())
		assert.Nil(t, v)
		assert.Equal(t, ErrEdgeNotFound, err)

		v, err = g.TargetVertexOf(NewEdge())
		assert.Nil(t, v)
		assert.Equal(t, ErrEdgeNotFound, err)
	}
}
