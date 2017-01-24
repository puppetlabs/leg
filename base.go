// Portions of this file are derived from JGraphT, a free Java graph-theory
// library.
//
// (C) Copyright 2003-2016, by Barak Naveh and Contributors.

package gographt

import (
	"github.com/reflect/godat"
)

type IntrusiveEdge struct {
	Source, Target Vertex
	Edge           Edge
	Weight         float64
}

type BaseGraphOps interface {
	EdgesBetween(source, target Vertex) EdgeSet
	EdgeBetween(source, target Vertex) (Edge, error)
	EdgesOf(vertex Vertex) EdgeSet
	AddEdge(edge Edge)
	RemoveEdge(edge Edge)
	Vertices() MutableVertexSet
}

type baseEdgesView struct {
	g *BaseGraph
}

func (sev *baseEdgesView) Contains(edge Edge) bool {
	return sev.g.edges.Contains(edge)
}

func (sev *baseEdgesView) Count() uint {
	return uint(sev.g.edges.Size())
}

func (sev *baseEdgesView) AsSlice() []Edge {
	s := make([]Edge, sev.g.edges.Size())

	i := 0
	sev.ForEach(func(edge Edge) error {
		s[i] = edge
		i++

		return nil
	})

	return s
}

func (sev *baseEdgesView) ForEach(fn EdgeSetIterationFunc) error {
	return sev.g.edges.ForEachInto(func(key Edge, value *IntrusiveEdge) error {
		return fn(key)
	})
}

type BaseGraph struct {
	AllowsLoops, AllowsMultipleEdges bool
	Ops                              BaseGraphOps

	features  GraphFeature
	edges     godat.Map // map[Edge]*IntrusiveEdge
	edgesView EdgeSet
}

func (g *BaseGraph) Features() GraphFeature {
	return g.features
}

func (g *BaseGraph) EdgesBetween(source, target Vertex) EdgeSet {
	return g.Ops.EdgesBetween(source, target)
}

func (g *BaseGraph) EdgeBetween(source, target Vertex) (Edge, error) {
	return g.Ops.EdgeBetween(source, target)
}

func (g *BaseGraph) Connect(source, target Vertex) error {
	return g.AddEdge(source, target, NewEdge())
}

func (g *BaseGraph) AddEdge(source, target Vertex, edge Edge) error {
	return g.AddEdgeWithWeight(source, target, edge, DefaultEdgeWeight)
}

func (g *BaseGraph) ConnectWithWeight(source, target Vertex, weight float64) error {
	return g.AddEdgeWithWeight(source, target, NewEdge(), weight)
}

func (g *BaseGraph) AddEdgeWithWeight(source, target Vertex, edge Edge, weight float64) error {
	if g.ContainsEdge(edge) {
		return ErrEdgeAlreadyInGraph
	}

	if !g.ContainsVertex(source) {
		return &VertexNotFoundError{source}
	}
	if !g.ContainsVertex(target) {
		return &VertexNotFoundError{target}
	}

	if !g.AllowsMultipleEdges && g.ContainsEdgeBetween(source, target) {
		return ErrEdgeAlreadyInGraph
	}

	if !g.AllowsLoops && source == target {
		return ErrWouldCreateLoop
	}

	ie := &IntrusiveEdge{
		Source: source,
		Target: target,
		Edge:   edge,
		Weight: weight,
	}

	g.edges.Put(edge, ie)
	g.Ops.AddEdge(edge)

	return nil
}

func (g *BaseGraph) AddVertex(vertex Vertex) {
	g.Ops.Vertices().Add(vertex)
}

func (g *BaseGraph) ContainsEdgeBetween(source, target Vertex) bool {
	_, err := g.EdgeBetween(source, target)
	return err == nil
}

func (g *BaseGraph) ContainsEdge(edge Edge) bool {
	return g.edges.Contains(edge)
}

func (g *BaseGraph) ContainsVertex(vertex Vertex) bool {
	return g.Vertices().Contains(vertex)
}

func (g *BaseGraph) Edges() EdgeSet {
	if g.edgesView == nil {
		g.edgesView = &baseEdgesView{g}
	}

	return g.edgesView
}

func (g *BaseGraph) EdgesOf(vertex Vertex) (EdgeSet, error) {
	if !g.ContainsVertex(vertex) {
		return nil, &VertexNotFoundError{vertex}
	}

	return g.Ops.EdgesOf(vertex), nil
}

func (g *BaseGraph) RemoveEdges(edges []Edge) (modified bool) {
	for _, edge := range edges {
		modified = modified || g.RemoveEdge(edge)
	}

	return
}

func (g *BaseGraph) RemoveEdgesBetween(source, target Vertex) EdgeSet {
	edges := g.EdgesBetween(source, target)
	g.RemoveEdges(edges.AsSlice())

	return edges
}

func (g *BaseGraph) RemoveEdge(edge Edge) bool {
	if !g.ContainsEdge(edge) {
		return false
	}

	g.Ops.RemoveEdge(edge)
	g.edges.Remove(edge)

	return true
}

func (g *BaseGraph) RemoveEdgeBetween(source, target Vertex) (Edge, error) {
	edge, err := g.EdgeBetween(source, target)
	if err != nil {
		return nil, err
	}

	g.RemoveEdge(edge)
	return edge, nil
}

func (g *BaseGraph) RemoveVertices(vertices []Vertex) (modified bool) {
	for _, vertex := range vertices {
		modified = modified || g.RemoveVertex(vertex)
	}

	return
}

func (g *BaseGraph) RemoveVertex(vertex Vertex) bool {
	if !g.ContainsVertex(vertex) {
		return false
	}

	g.RemoveEdges(g.Ops.EdgesOf(vertex).AsSlice())
	g.Ops.Vertices().Remove(vertex)

	return true
}

func (g *BaseGraph) Vertices() VertexSet {
	return g.Ops.Vertices()
}

func (g *BaseGraph) SourceVertexOf(edge Edge) (Vertex, error) {
	ie, found := g.edges.Get(edge)
	if !found {
		return nil, ErrEdgeNotFound
	}

	return ie.(*IntrusiveEdge).Source, nil
}

func (g *BaseGraph) TargetVertexOf(edge Edge) (Vertex, error) {
	ie, found := g.edges.Get(edge)
	if !found {
		return nil, ErrEdgeNotFound
	}

	return ie.(*IntrusiveEdge).Target, nil
}

func (g *BaseGraph) WeightOf(edge Edge) (float64, error) {
	ie, found := g.edges.Get(edge)
	if !found {
		return DefaultEdgeWeight, ErrEdgeNotFound
	}

	return ie.(*IntrusiveEdge).Weight, nil
}

func NewBaseGraph(features GraphFeature, allowsLoops, allowsMultipleEdges bool, ops BaseGraphOps) *BaseGraph {
	var edges godat.Map
	if features&DeterministicIteration != 0 {
		edges = godat.NewLinkedHashMap()
	} else {
		edges = godat.NewHashMap()
	}

	return &BaseGraph{
		AllowsLoops:         allowsLoops,
		AllowsMultipleEdges: allowsMultipleEdges,
		Ops:                 ops,

		features: features,
		edges:    edges,
	}
}

type BaseUndirectedGraph struct {
	*BaseGraph
	Ops BaseUndirectedGraphOps
}

type BaseUndirectedGraphOps interface {
	BaseGraphOps
	DegreeOf(vertex Vertex) uint
}

func (ug *BaseUndirectedGraph) DegreeOf(vertex Vertex) (uint, error) {
	if !ug.ContainsVertex(vertex) {
		return 0, &VertexNotFoundError{vertex}
	}

	return ug.Ops.DegreeOf(vertex), nil
}

func NewBaseUndirectedGraph(features GraphFeature, allowsLoops, allowsMultipleEdges bool, ops BaseUndirectedGraphOps) *BaseUndirectedGraph {
	return &BaseUndirectedGraph{NewBaseGraph(features, allowsLoops, allowsMultipleEdges, ops), ops}
}

type BaseDirectedGraph struct {
	*BaseGraph
	Ops BaseDirectedGraphOps
}

type BaseDirectedGraphOps interface {
	BaseGraphOps
	InDegreeOf(vertex Vertex) uint
	IncomingEdgesOf(vertex Vertex) EdgeSet
	OutDegreeOf(vertex Vertex) uint
	OutgoingEdgesOf(vertex Vertex) EdgeSet
}

func (dg *BaseDirectedGraph) InDegreeOf(vertex Vertex) (uint, error) {
	if !dg.ContainsVertex(vertex) {
		return 0, &VertexNotFoundError{vertex}
	}

	return dg.Ops.InDegreeOf(vertex), nil
}

func (dg *BaseDirectedGraph) IncomingEdgesOf(vertex Vertex) (EdgeSet, error) {
	if !dg.ContainsVertex(vertex) {
		return nil, &VertexNotFoundError{vertex}
	}

	return dg.Ops.IncomingEdgesOf(vertex), nil
}

func (dg *BaseDirectedGraph) OutDegreeOf(vertex Vertex) (uint, error) {
	if !dg.ContainsVertex(vertex) {
		return 0, &VertexNotFoundError{vertex}
	}

	return dg.Ops.OutDegreeOf(vertex), nil
}

func (dg *BaseDirectedGraph) OutgoingEdgesOf(vertex Vertex) (EdgeSet, error) {
	if !dg.ContainsVertex(vertex) {
		return nil, &VertexNotFoundError{vertex}
	}

	return dg.Ops.OutgoingEdgesOf(vertex), nil
}

func NewBaseDirectedGraph(features GraphFeature, allowsLoops, allowsMultipleEdges bool, ops BaseDirectedGraphOps) *BaseDirectedGraph {
	return &BaseDirectedGraph{NewBaseGraph(features, allowsLoops, allowsMultipleEdges, ops), ops}
}
