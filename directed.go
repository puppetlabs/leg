// Portions of this file are derived from JGraphT, a free Java graph-theory
// library.
//
// (C) Copyright 2015-2017, by Barak Naveh and Contributors.

package gographt

import (
	"github.com/reflect/godat"
)

type directedEdgeContainer struct {
	incoming MutableEdgeSet
	outgoing MutableEdgeSet
}

type directedVertexSet struct {
	features GraphFeature
	storage  godat.Map // map[Vertex]*directedEdgeContainer
}

func (vs *directedVertexSet) Contains(vertex Vertex) bool {
	return vs.storage.Contains(vertex)
}

func (vs *directedVertexSet) Count() uint {
	return uint(vs.storage.Size())
}

func (vs *directedVertexSet) AsSlice() []Vertex {
	s := make([]Vertex, 0, vs.Count())
	vs.storage.KeysInto(&s)
	return s
}

func (vs *directedVertexSet) ForEach(fn VertexSetIterationFunc) error {
	return vs.storage.ForEachInto(func(key Vertex, value *directedEdgeContainer) error {
		return fn(key)
	})
}

func (vs *directedVertexSet) Add(vertex Vertex) {
	if vs.storage.Contains(vertex) {
		return
	}

	vs.storage.Put(vertex, nil)
}

func (vs *directedVertexSet) Remove(vertex Vertex) {
	vs.storage.Remove(vertex)
}

func (vs *directedVertexSet) initVertex(vertex Vertex) *directedEdgeContainer {
	if !vs.storage.Contains(vertex) {
		return nil
	}

	var container *directedEdgeContainer
	vs.storage.GetInto(vertex, &container)

	if container == nil {
		container = &directedEdgeContainer{}
		vs.storage.Put(vertex, container)
	}

	return container
}

func (vs *directedVertexSet) incomingEdgesOf(vertex Vertex) MutableEdgeSet {
	container := vs.initVertex(vertex)
	if container == nil {
		return nil
	}

	if container.incoming == nil {
		if vs.features&DeterministicIteration != 0 {
			container.incoming = NewMutableEdgeSet(godat.NewLinkedHashSet())
		} else {
			container.incoming = NewMutableEdgeSet(godat.NewHashSet())
		}
	}

	return container.incoming
}

func (vs *directedVertexSet) outgoingEdgesOf(vertex Vertex) MutableEdgeSet {
	container := vs.initVertex(vertex)
	if container == nil {
		return nil
	}

	if container.outgoing == nil {
		if vs.features&DeterministicIteration != 0 {
			container.outgoing = NewMutableEdgeSet(godat.NewLinkedHashSet())
		} else {
			container.outgoing = NewMutableEdgeSet(godat.NewHashSet())
		}
	}

	return container.outgoing
}

type directedGraphOps struct {
	g        *baseDirectedGraph
	vertices *directedVertexSet
}

func (o *directedGraphOps) EdgesBetween(source, target Vertex) EdgeSet {
	if !o.g.ContainsVertex(source) || !o.g.ContainsVertex(target) {
		return nil
	}

	es := &unenforcedSliceEdgeSet{}

	o.vertices.outgoingEdgesOf(source).ForEach(func(edge Edge) error {
		tt, _ := o.g.TargetVertexOf(edge)
		if tt == target {
			es.Add(edge)
		}

		return nil
	})

	return es
}

func (o *directedGraphOps) EdgeBetween(source, target Vertex) (Edge, error) {
	if !o.g.ContainsVertex(source) || !o.g.ContainsVertex(target) {
		return nil, &NotConnectedError{Source: source, Target: target}
	}

	var found Edge
	err := o.vertices.outgoingEdgesOf(source).ForEach(func(edge Edge) error {
		tt, _ := o.g.TargetVertexOf(edge)
		if tt == target {
			found = edge
			return godat.ErrStopIteration
		}

		return nil
	})

	if err == godat.ErrStopIteration {
		return found, nil
	}

	return nil, &NotConnectedError{Source: source, Target: target}
}

func (o *directedGraphOps) EdgesOf(vertex Vertex) EdgeSet {
	if !o.g.ContainsVertex(vertex) {
		return nil
	}

	var set MutableEdgeSet
	if o.g.Features()&DeterministicIteration != 0 {
		set = NewMutableEdgeSet(godat.NewLinkedHashSet())
	} else {
		set = NewMutableEdgeSet(godat.NewHashSet())
	}

	o.IncomingEdgesOf(vertex).ForEach(func(edge Edge) error {
		set.Add(edge)
		return nil
	})
	o.OutgoingEdgesOf(vertex).ForEach(func(edge Edge) error {
		set.Add(edge)
		return nil
	})

	return set
}

func (o *directedGraphOps) AddEdge(edge Edge) {
	source, _ := o.g.SourceVertexOf(edge)
	target, _ := o.g.TargetVertexOf(edge)

	o.vertices.outgoingEdgesOf(source).Add(edge)
	o.vertices.incomingEdgesOf(target).Add(edge)
}

func (o *directedGraphOps) RemoveEdge(edge Edge) {
	source, _ := o.g.SourceVertexOf(edge)
	target, _ := o.g.TargetVertexOf(edge)

	o.vertices.outgoingEdgesOf(source).Remove(edge)
	o.vertices.incomingEdgesOf(target).Remove(edge)
}

func (o *directedGraphOps) Vertices() MutableVertexSet {
	return o.vertices
}

func (o *directedGraphOps) InDegreeOf(vertex Vertex) uint {
	if !o.g.ContainsVertex(vertex) {
		return 0
	}

	return o.vertices.incomingEdgesOf(vertex).Count()
}

func (o *directedGraphOps) IncomingEdgesOf(vertex Vertex) EdgeSet {
	if !o.g.ContainsVertex(vertex) {
		return nil
	}

	return o.vertices.incomingEdgesOf(vertex)
}

func (o *directedGraphOps) OutDegreeOf(vertex Vertex) uint {
	if !o.g.ContainsVertex(vertex) {
		return 0
	}

	return o.vertices.outgoingEdgesOf(vertex).Count()
}

func (o *directedGraphOps) OutgoingEdgesOf(vertex Vertex) EdgeSet {
	if !o.g.ContainsVertex(vertex) {
		return nil
	}

	return o.vertices.outgoingEdgesOf(vertex)
}

func newDirectedGraph(features GraphFeature, allowLoops, allowMultipleEdges bool) *baseDirectedGraph {
	var vertexStorage godat.Map
	if features&DeterministicIteration != 0 {
		vertexStorage = godat.NewLinkedHashMap()
	} else {
		vertexStorage = godat.NewHashMap()
	}

	ops := &directedGraphOps{
		vertices: &directedVertexSet{features: features, storage: vertexStorage},
	}

	g := newBaseDirectedGraph(features, allowLoops, allowMultipleEdges, ops)
	ops.g = g

	return g
}

//
// Simple graphs
//

type SimpleDirectedGraph struct {
	DirectedGraph
}

// Creates a new simple graph.
//
// A simple graph is an directed graph that does not permit loops or multiple
// edges between vertices.
func NewSimpleDirectedGraph() *SimpleDirectedGraph {
	return NewSimpleDirectedGraphWithFeatures(0)
}

func NewSimpleDirectedGraphWithFeatures(features GraphFeature) *SimpleDirectedGraph {
	return &SimpleDirectedGraph{newDirectedGraph(features, false, false)}
}

type SimpleDirectedWeightedGraph struct {
	DirectedWeightedGraph
}

// Creates a new simple weighted graph.
//
// A simple weighted graph is a simple graph for which edges are assigned
// weights.
func NewSimpleDirectedWeightedGraph() *SimpleDirectedWeightedGraph {
	return NewSimpleDirectedWeightedGraphWithFeatures(0)
}

func NewSimpleDirectedWeightedGraphWithFeatures(features GraphFeature) *SimpleDirectedWeightedGraph {
	return &SimpleDirectedWeightedGraph{DirectedWeightedGraph(newDirectedGraph(features, false, false))}
}

type DirectedMultigraph struct {
	DirectedGraph
}

//
// Multigraphs
//

// Creates a new multigraph.
//
// A multigraph is an directed graph that does not permit loops, but does
// permit multiple edges between any two vertices.
func NewDirectedMultigraph() *DirectedMultigraph {
	return NewDirectedMultigraphWithFeatures(0)
}

func NewDirectedMultigraphWithFeatures(features GraphFeature) *DirectedMultigraph {
	return &DirectedMultigraph{newDirectedGraph(features, false, true)}
}

type DirectedWeightedMultigraph struct {
	DirectedWeightedGraph
}

// Creates a new weighted multigraph.
func NewDirectedWeightedMultigraph() *DirectedWeightedMultigraph {
	return NewDirectedWeightedMultigraphWithFeatures(0)
}

func NewDirectedWeightedMultigraphWithFeatures(features GraphFeature) *DirectedWeightedMultigraph {
	return &DirectedWeightedMultigraph{DirectedWeightedGraph(newDirectedGraph(features, false, true))}
}

//
// Pseudographs
//

type DirectedPseudograph struct {
	DirectedGraph
}

// Creates a new pseudograph.
//
// A pseudograph is an directed graph that permits both loops and multiple
// edges between vertices.
func NewDirectedPseudograph() *DirectedPseudograph {
	return NewDirectedPseudographWithFeatures(0)
}

func NewDirectedPseudographWithFeatures(features GraphFeature) *DirectedPseudograph {
	return &DirectedPseudograph{newDirectedGraph(features, true, true)}
}

type DirectedWeightedPseudograph struct {
	DirectedWeightedGraph
}

// Creates a new weighted pseudograph.
func NewDirectedWeightedPseudograph() *DirectedWeightedPseudograph {
	return NewDirectedWeightedPseudographWithFeatures(0)
}

func NewDirectedWeightedPseudographWithFeatures(features GraphFeature) *DirectedWeightedPseudograph {
	return &DirectedWeightedPseudograph{DirectedWeightedGraph(newDirectedGraph(features, true, true))}
}
