// Portions of this file are derived from JGraphT, a free Java graph-theory
// library.
//
// (C) Copyright 2003-2016, by Barak Naveh and Contributors.

package gographt

type Vertex interface{}

type VertexSetIterationFunc func(vertex Vertex) error

type VertexSet interface {
	Contains(vertex Vertex) bool
	Count() uint
	AsSlice() []Vertex
	ForEach(fn VertexSetIterationFunc) error
}

type MutableVertexSet interface {
	VertexSet

	Add(vertex Vertex)
	Remove(vertex Vertex)
}

type Edge interface{}

type EdgeSetIterationFunc func(edge Edge) error

type EdgeSet interface {
	Contains(edge Edge) bool
	Count() uint
	AsSlice() []Edge
	ForEach(fn EdgeSetIterationFunc) error
}

type MutableEdgeSet interface {
	EdgeSet

	Add(edge Edge)
	Remove(edge Edge)
}

const DefaultEdgeWeight = float64(1.)

type Graph interface {
	Features() GraphFeature
	EdgesBetween(source, target Vertex) EdgeSet
	EdgeBetween(source, target Vertex) (Edge, error)
	Connect(source, target Vertex) error
	AddEdge(source, target Vertex, edge Edge) error
	AddVertex(vertex Vertex)
	ContainsEdgeBetween(source, target Vertex) bool
	ContainsEdge(edge Edge) bool
	ContainsVertex(vertex Vertex) bool
	Edges() EdgeSet
	EdgesOf(vertex Vertex) (EdgeSet, error)
	RemoveEdges(edges []Edge) bool
	RemoveEdgesBetween(source, target Vertex) EdgeSet
	RemoveEdge(edge Edge) bool
	RemoveEdgeBetween(source, target Vertex) (Edge, error)
	RemoveVertices(vertices []Vertex) bool
	RemoveVertex(vertex Vertex) bool
	Vertices() VertexSet
	SourceVertexOf(edge Edge) (Vertex, error)
	TargetVertexOf(edge Edge) (Vertex, error)
	WeightOf(edge Edge) (float64, error)
}

type DirectedGraph interface {
	Graph

	InDegreeOf(vertex Vertex) (uint, error)
	IncomingEdgesOf(vertex Vertex) (EdgeSet, error)
	OutDegreeOf(vertex Vertex) (uint, error)
	OutgoingEdgesOf(vertex Vertex) (EdgeSet, error)
}

type UndirectedGraph interface {
	Graph

	DegreeOf(vertex Vertex) (uint, error)
}

type WeightedGraph interface {
	ConnectWithWeight(source, target Vertex, weight float64) error
	AddEdgeWithWeight(source, target Vertex, edge Edge, weight float64) error
}

type DirectedWeightedGraph interface {
	DirectedGraph
	WeightedGraph
}

type UndirectedWeightedGraph interface {
	UndirectedGraph
	WeightedGraph
}
