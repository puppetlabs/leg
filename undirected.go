// Portions of this file are derived from JGraphT, a free Java graph-theory
// library.
//
// (C) Copyright 2015-2016, by Barak Naveh and Contributors.

package gographt

import (
	"github.com/reflect/gographt/data"
)

type undirectedVertexSet struct {
	features GraphFeature
	storage  data.Map // map[Vertex]MutableEdgeSet
}

func (vs *undirectedVertexSet) Contains(vertex Vertex) bool {
	return vs.storage.Contains(vertex)
}

func (vs *undirectedVertexSet) Count() uint {
	return uint(vs.storage.Size())
}

func (vs *undirectedVertexSet) AsSlice() []Vertex {
	s := make([]Vertex, vs.Count())

	i := 0
	vs.ForEach(func(vertex Vertex) error {
		s[i] = vertex
		i++

		return nil
	})

	return s
}

func (vs *undirectedVertexSet) ForEach(fn VertexSetIterationFunc) error {
	return vs.storage.ForEach(func(key, value interface{}) error {
		return fn(key.(Vertex))
	})
}

func (vs *undirectedVertexSet) Add(vertex Vertex) {
	if vs.storage.Contains(vertex) {
		return
	}

	vs.storage.Put(vertex, nil)
}

func (vs *undirectedVertexSet) Remove(vertex Vertex) {
	vs.storage.Remove(vertex)
}

func (vs *undirectedVertexSet) edgesOf(vertex Vertex) MutableEdgeSet {
	if !vs.storage.Contains(vertex) {
		return nil
	}

	var set MutableEdgeSet
	if s, _ := vs.storage.Get(vertex); s != nil {
		set = s.(MutableEdgeSet)
	} else {
		if vs.features&DeterministicIteration != 0 {
			set = NewMutableEdgeSet(data.NewLinkedHashSet())
		} else {
			set = NewMutableEdgeSet(data.NewHashSet())
		}

		vs.storage.Put(vertex, set)
	}

	return set
}

type undirectedGraphOps struct {
	g        *BaseUndirectedGraph
	vertices *undirectedVertexSet
}

func (o *undirectedGraphOps) EdgesBetween(source, target Vertex) EdgeSet {
	if !o.g.ContainsVertex(source) || !o.g.ContainsVertex(target) {
		return nil
	}

	es := &UnenforcedSliceEdgeSet{}

	o.vertices.edgesOf(source).ForEach(func(edge Edge) error {
		if o.edgeHasSourceAndTarget(edge, source, target) {
			es.Add(edge)
		}

		return nil
	})

	return es
}

func (o *undirectedGraphOps) EdgeBetween(source, target Vertex) (Edge, error) {
	if !o.g.ContainsVertex(source) || !o.g.ContainsVertex(target) {
		return nil, &NotConnectedError{Source: source, Target: target}
	}

	var found Edge
	err := o.vertices.edgesOf(source).ForEach(func(edge Edge) error {
		if o.edgeHasSourceAndTarget(edge, source, target) {
			found = edge
			return data.ErrStopIteration
		}

		return nil
	})

	if err == data.ErrStopIteration {
		return found, nil
	}

	return nil, &NotConnectedError{Source: source, Target: target}
}

func (o *undirectedGraphOps) edgeHasSourceAndTarget(edge Edge, source, target Vertex) bool {
	ts, _ := o.g.SourceVertexOf(edge)
	tt, _ := o.g.TargetVertexOf(edge)

	return (source == ts && target == tt) || (source == tt && target == ts)
}

func (o *undirectedGraphOps) EdgesOf(vertex Vertex) EdgeSet {
	if !o.g.ContainsVertex(vertex) {
		return nil
	}

	return o.vertices.edgesOf(vertex)
}

func (o *undirectedGraphOps) AddEdge(edge Edge) {
	source, _ := o.g.SourceVertexOf(edge)
	target, _ := o.g.TargetVertexOf(edge)

	o.vertices.edgesOf(source).Add(edge)
	o.vertices.edgesOf(target).Add(edge)
}

func (o *undirectedGraphOps) RemoveEdge(edge Edge) {
	source, _ := o.g.SourceVertexOf(edge)
	target, _ := o.g.TargetVertexOf(edge)

	o.vertices.edgesOf(source).Remove(edge)
	o.vertices.edgesOf(target).Remove(edge)
}

func (o *undirectedGraphOps) Vertices() MutableVertexSet {
	return o.vertices
}

func (o *undirectedGraphOps) DegreeOf(vertex Vertex) uint {
	if !o.g.ContainsVertex(vertex) {
		return 0
	}

	return o.vertices.edgesOf(vertex).Count()
}

func newUndirectedGraph(features GraphFeature, allowLoops, allowMultipleEdges bool) *BaseUndirectedGraph {
	var vertexStorage data.Map
	if features&DeterministicIteration != 0 {
		vertexStorage = data.NewLinkedHashMap()
	} else {
		vertexStorage = data.NewHashMap()
	}

	ops := &undirectedGraphOps{
		vertices: &undirectedVertexSet{features: features, storage: vertexStorage},
	}

	g := NewBaseUndirectedGraph(features, allowLoops, allowMultipleEdges, ops)
	ops.g = g

	return g
}

//
// Simple graphs
//

type SimpleGraph struct {
	UndirectedGraph
}

// Creates a new simple graph.
//
// A simple graph is an undirected graph that does not permit loops or multiple
// edges between vertices.
func NewSimpleGraph() *SimpleGraph {
	return NewSimpleGraphWithFeatures(0)
}

func NewSimpleGraphWithFeatures(features GraphFeature) *SimpleGraph {
	return &SimpleGraph{newUndirectedGraph(features, false, false)}
}

type SimpleWeightedGraph struct {
	UndirectedWeightedGraph
}

// Creates a new simple weighted graph.
//
// A simple weighted graph is a simple graph for which edges are assigned
// weights.
func NewSimpleWeightedGraph() *SimpleWeightedGraph {
	return NewSimpleWeightedGraphWithFeatures(0)
}

func NewSimpleWeightedGraphWithFeatures(features GraphFeature) *SimpleWeightedGraph {
	return &SimpleWeightedGraph{UndirectedWeightedGraph(newUndirectedGraph(features, false, false))}
}

type UndirectedMultigraph struct {
	UndirectedGraph
}

//
// Multigraphs
//

// Creates a new multigraph.
//
// A multigraph is an undirected graph that does not permit loops, but does
// permit multiple edges between any two vertices.
func NewUndirectedMultigraph() *UndirectedMultigraph {
	return NewUndirectedMultigraphWithFeatures(0)
}

func NewUndirectedMultigraphWithFeatures(features GraphFeature) *UndirectedMultigraph {
	return &UndirectedMultigraph{newUndirectedGraph(features, false, true)}
}

type UndirectedWeightedMultigraph struct {
	UndirectedWeightedGraph
}

// Creates a new weighted multigraph.
func NewUndirectedWeightedMultigraph() *UndirectedWeightedMultigraph {
	return NewUndirectedWeightedMultigraphWithFeatures(0)
}

func NewUndirectedWeightedMultigraphWithFeatures(features GraphFeature) *UndirectedWeightedMultigraph {
	return &UndirectedWeightedMultigraph{UndirectedWeightedGraph(newUndirectedGraph(features, false, true))}
}

//
// Pseudographs
//

type UndirectedPseudograph struct {
	UndirectedGraph
}

// Creates a new pseudograph.
//
// A pseudograph is an undirected graph that permits both loops and multiple
// edges between vertices.
func NewUndirectedPseudograph() *UndirectedPseudograph {
	return NewUndirectedPseudographWithFeatures(0)
}

func NewUndirectedPseudographWithFeatures(features GraphFeature) *UndirectedPseudograph {
	return &UndirectedPseudograph{newUndirectedGraph(features, true, true)}
}

type UndirectedWeightedPseudograph struct {
	UndirectedWeightedGraph
}

// Creates a new weighted pseudograph.
func NewUndirectedWeightedPseudograph() *UndirectedWeightedPseudograph {
	return NewUndirectedWeightedPseudographWithFeatures(0)
}

func NewUndirectedWeightedPseudographWithFeatures(features GraphFeature) *UndirectedWeightedPseudograph {
	return &UndirectedWeightedPseudograph{UndirectedWeightedGraph(newUndirectedGraph(features, true, true))}
}
