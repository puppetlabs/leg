// Portions of this file are derived from JGraphT, a free Java graph-theory
// library.
//
// (C) Copyright 2013-2016, by Alexey Kudinkin and Contributors.

package algo

import (
	"github.com/reflect/gographt"
	"github.com/reflect/gographt/data"
)

type PrimMinimumSpanningTree struct {
	TotalWeight float64

	es gographt.MutableEdgeSet
}

func (mst *PrimMinimumSpanningTree) Edges() gographt.EdgeSet {
	return mst.es
}

func PrimMinimumSpanningTreeOf(g gographt.UndirectedWeightedGraph) *PrimMinimumSpanningTree {
	vs := g.Vertices()

	var (
		es        gographt.MutableEdgeSet
		unspanned data.Set
	)

	if g.Features()&gographt.DeterministicIteration != 0 {
		es = gographt.NewMutableEdgeSet(data.NewLinkedHashSet())
		unspanned = data.NewLinkedHashSetWithCapacity(int(vs.Count()))
	} else {
		es = gographt.NewMutableEdgeSet(data.NewHashSet())
		unspanned = data.NewHashSetWithCapacity(int(vs.Count()))
	}

	mst := &PrimMinimumSpanningTree{es: es}

	vs.ForEach(func(vertex gographt.Vertex) error {
		unspanned.Add(vertex)
		return nil
	})

	for !unspanned.Empty() {
		var root gographt.Vertex
		unspanned.ForEach(func(element interface{}) error {
			root = element.(gographt.Vertex)
			return data.ErrStopIteration
		})

		unspanned.Remove(root)

		dangling := data.NewPriorityQueue()

		edges, _ := g.EdgesOf(root)
		edges.ForEach(func(edge gographt.Edge) error {
			weight, _ := g.WeightOf(edge)
			dangling.Add(edge, -weight)

			return nil
		})

		var next gographt.Edge
		for dangling.Poll(&next) {
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
