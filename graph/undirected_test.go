package graph

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	SimpleGraphConstructors = []func() MutableUndirectedGraph{
		func() MutableUndirectedGraph { return NewSimpleGraph() },
		func() MutableUndirectedGraph { return NewSimpleGraphWithFeatures(DeterministicIteration) },
		func() MutableUndirectedGraph { return NewSimpleWeightedGraph() },
		func() MutableUndirectedGraph { return NewSimpleWeightedGraphWithFeatures(DeterministicIteration) },
	}

	UndirectedMultigraphConstructors = []func() MutableUndirectedGraph{
		func() MutableUndirectedGraph { return NewUndirectedMultigraph() },
		func() MutableUndirectedGraph { return NewUndirectedMultigraphWithFeatures(DeterministicIteration) },
		func() MutableUndirectedGraph { return NewUndirectedWeightedMultigraph() },
		func() MutableUndirectedGraph {
			return NewUndirectedWeightedMultigraphWithFeatures(DeterministicIteration)
		},
	}

	UndirectedPseudographConstructors = []func() MutableUndirectedGraph{
		func() MutableUndirectedGraph { return NewUndirectedPseudograph() },
		func() MutableUndirectedGraph { return NewUndirectedPseudographWithFeatures(DeterministicIteration) },
		func() MutableUndirectedGraph { return NewUndirectedWeightedPseudograph() },
		func() MutableUndirectedGraph {
			return NewUndirectedWeightedPseudographWithFeatures(DeterministicIteration)
		},
	}

	UndirectedConstructors = append(
		append(
			append([]func() MutableUndirectedGraph{}, SimpleGraphConstructors...),
			UndirectedMultigraphConstructors...),
		UndirectedPseudographConstructors...)
)

func TestUndirectedMultigraphEdgeAddition(t *testing.T) {
	for _, constructor := range UndirectedMultigraphConstructors {
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

		assert.NoError(t, g.AddEdge("a", "b", 2))

		edges := g.EdgesBetween("a", "b")
		assert.Equal(t, uint(2), edges.Count())
		assert.True(t, edges.Contains(1))
		assert.True(t, edges.Contains(2))
	}
}

func TestUndirectedConnections(t *testing.T) {
	for _, constructor := range UndirectedConstructors {
		g := constructor()
		g.AddVertex("a")
		g.AddVertex("b")
		g.AddVertex("c")
		g.AddVertex("d")

		assert.NoError(t, g.AddEdge("a", "b", 1))
		assert.NoError(t, g.AddEdge("b", "c", 2))
		assert.NoError(t, g.AddEdge("c", "d", 3))

		counts := map[string]uint{
			"a": 1,
			"b": 2,
			"c": 2,
			"d": 1,
		}

		for vertex, expected := range counts {
			n, err := g.DegreeOf(vertex)
			assert.NoError(t, err)
			assert.Equal(t, expected, n, "for vertex %s", vertex)
		}

		assert.True(t, g.RemoveEdge(1))

		// We can't remove the same edge again, obviously.
		assert.False(t, g.RemoveEdge(1))

		counts["a"]--
		counts["b"]--

		for vertex, expected := range counts {
			n, err := g.DegreeOf(vertex)
			assert.NoError(t, err)
			assert.Equal(t, expected, n)
		}

		edge, err := g.RemoveEdgeBetween("b", "c")
		assert.NoError(t, err)
		assert.Equal(t, 2, edge)

		counts["b"]--
		counts["c"]--

		for vertex, expected := range counts {
			n, err := g.DegreeOf(vertex)
			assert.NoError(t, err)
			assert.Equal(t, expected, n, "for vertex %s", vertex)
		}

		assert.True(t, g.RemoveVertex("c"))

		// We can't remove the same vertex either.
		assert.False(t, g.RemoveVertex("c"))

		delete(counts, "c")
		counts["d"]--

		for vertex, expected := range counts {
			n, err := g.DegreeOf(vertex)
			assert.NoError(t, err)
			assert.Equal(t, expected, n, "for vertex %s", vertex)
		}

		_, err = g.DegreeOf("c")
		assert.Equal(t, &VertexNotFoundError{Vertex: "c"}, err)
	}
}

func TestUndirectedVertexNotFound(t *testing.T) {
	for _, constructor := range UndirectedConstructors {
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

func TestUndirectedEdgeNotFound(t *testing.T) {
	for _, constructor := range UndirectedConstructors {
		g := constructor()
		g.AddVertex("a")
		g.AddVertex("b")
		g.AddVertex("c")
		g.Connect("b", "c")

		edge, err := g.EdgeBetween("a", "b")
		assert.Nil(t, edge)
		assert.Equal(t, &NotConnectedError{Source: "a", Target: "b"}, err)

		v, err := g.SourceVertexOf(NewEdge())
		assert.Nil(t, v)
		assert.Equal(t, ErrEdgeNotFound, err)

		v, err = g.TargetVertexOf(NewEdge())
		assert.Nil(t, v)
		assert.Equal(t, ErrEdgeNotFound, err)
	}
}

func ExampleSimpleGraph() {
	g := NewSimpleGraph()
	g.AddVertex("Charles")
	g.AddVertex("Manson")
	g.AddVertex("Marilyn")
	g.AddVertex("Monroe")

	g.Connect("Charles", "Manson")
	g.Connect("Marilyn", "Monroe")
	g.Connect("Marilyn", "Manson")

	fmt.Println("Vertices:", g.Vertices().Count())
	fmt.Println("Edges:", g.Edges().Count())
	// Output:
	// Vertices: 4
	// Edges: 3
}
