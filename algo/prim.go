// Portions of this file are derived from JGraphT, a free Java graph-theory
// library.
//
// (C) Copyright 2013-2016, by Alexey Kudinkin and Contributors.

package algo

import (
	"github.com/reflect/godat"
	"github.com/reflect/gographt"
)

const (
	PrimMinimumSpanningTreeSupportedFeatures = gographt.DeterministicIteration
)

type PrimMinimumSpanningTree struct {
	TotalWeight float64

	features gographt.GraphFeature
	es       gographt.MutableEdgeSet
}

func (mst *PrimMinimumSpanningTree) Features() gographt.GraphFeature {
	return mst.features
}

func (mst *PrimMinimumSpanningTree) Edges() gographt.EdgeSet {
	return mst.es
}

func PrimMinimumSpanningTreeOf(g gographt.UndirectedWeightedGraph) *PrimMinimumSpanningTree {
	vs := g.Vertices()

	var (
		es        gographt.MutableEdgeSet
		unspanned godat.Set
	)

	if g.Features()&gographt.DeterministicIteration != 0 {
		es = gographt.NewMutableEdgeSet(godat.NewLinkedHashSet())
		unspanned = godat.NewLinkedHashSetWithCapacity(int(vs.Count()))
	} else {
		es = gographt.NewMutableEdgeSet(godat.NewHashSet())
		unspanned = godat.NewHashSetWithCapacity(int(vs.Count()))
	}

	mst := &PrimMinimumSpanningTree{
		features: g.Features() & PrimMinimumSpanningTreeSupportedFeatures,
		es:       es,
	}

	vs.ForEach(func(vertex gographt.Vertex) error {
		unspanned.Add(vertex)
		return nil
	})

	for !unspanned.Empty() {
		var root gographt.Vertex
		unspanned.ForEachInto(func(vertex gographt.Vertex) error {
			root = vertex
			return godat.ErrStopIteration
		})

		unspanned.Remove(root)

		dangling := godat.NewPriorityQueue()

		edges, _ := g.EdgesOf(root)
		edges.ForEach(func(edge gographt.Edge) error {
			weight, _ := g.WeightOf(edge)
			dangling.Add(edge, -weight)

			return nil
		})

		var next gographt.Edge
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
			edges.ForEach(func(edge gographt.Edge) error {
				candidate, _ := gographt.OppositeVertexOf(g, edge, target)
				if !unspanned.Contains(candidate) {
					return nil
				}

				weight, _ := g.WeightOf(edge)
				dangling.Add(edge, -weight)

				return nil
			})
		}
	}

	mst.es.ForEach(func(edge gographt.Edge) error {
		weight, _ := g.WeightOf(edge)
		mst.TotalWeight += weight

		return nil
	})

	return mst
}
