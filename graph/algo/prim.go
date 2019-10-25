// Portions of this file are derived from JGraphT, a free Java graph-theory
// library.
//
// (C) Copyright 2013-2016, by Alexey Kudinkin and Contributors.

package algo

import (
	"github.com/puppetlabs/horsehead/v2/datastructure"
	"github.com/puppetlabs/horsehead/v2/graph"
)

const (
	PrimMinimumSpanningTreeSupportedFeatures = graph.DeterministicIteration
)

type PrimMinimumSpanningTree struct {
	TotalWeight float64

	features graph.GraphFeature
	es       graph.MutableEdgeSet
}

func (mst *PrimMinimumSpanningTree) Features() graph.GraphFeature {
	return mst.features
}

func (mst *PrimMinimumSpanningTree) Edges() graph.EdgeSet {
	return mst.es
}

func PrimMinimumSpanningTreeOf(g graph.UndirectedGraph) *PrimMinimumSpanningTree {
	vs := g.Vertices()

	var (
		es        graph.MutableEdgeSet
		unspanned datastructure.Set
	)

	if g.Features()&graph.DeterministicIteration != 0 {
		es = graph.NewMutableEdgeSet(datastructure.NewLinkedHashSet())
		unspanned = datastructure.NewLinkedHashSetWithCapacity(int(vs.Count()))
	} else {
		es = graph.NewMutableEdgeSet(datastructure.NewHashSet())
		unspanned = datastructure.NewHashSetWithCapacity(int(vs.Count()))
	}

	mst := &PrimMinimumSpanningTree{
		features: g.Features() & PrimMinimumSpanningTreeSupportedFeatures,
		es:       es,
	}

	vs.ForEach(func(vertex graph.Vertex) error {
		unspanned.Add(vertex)
		return nil
	})

	for !unspanned.Empty() {
		var root graph.Vertex
		unspanned.ForEachInto(func(vertex graph.Vertex) error {
			root = vertex
			return datastructure.ErrStopIteration
		})

		unspanned.Remove(root)

		dangling := datastructure.NewPriorityQueue()

		edges, _ := g.EdgesOf(root)
		edges.ForEach(func(edge graph.Edge) error {
			weight, _ := g.WeightOf(edge)
			dangling.Add(edge, -weight)

			return nil
		})

		var next graph.Edge
		for dangling.PollInto(&next) {
			target, _ := g.SourceVertexOf(next)
			if !unspanned.Contains(target) {
				target, _ = g.TargetVertexOf(next)

				if !unspanned.Contains(target) {
					continue
				}
			}

			mst.es.Add(next)

			unspanned.Remove(target)

			edges, _ := g.EdgesOf(target)
			edges.ForEach(func(edge graph.Edge) error {
				candidate, _ := graph.OppositeVertexOf(g, edge, target)
				if !unspanned.Contains(candidate) {
					return nil
				}

				weight, _ := g.WeightOf(edge)
				dangling.Add(edge, -weight)

				return nil
			})
		}
	}

	mst.es.ForEach(func(edge graph.Edge) error {
		weight, _ := g.WeightOf(edge)
		mst.TotalWeight += weight

		return nil
	})

	return mst
}
