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

	var es gographt.MutableEdgeSet
	if g.Features()&gographt.DeterministicIteration != 0 {
		es = gographt.NewMutableEdgeSet(data.NewLinkedHashSet())
	} else {
		es = gographt.NewMutableEdgeSet(data.NewHashSet())
	}

	mst := &PrimMinimumSpanningTree{es: es}

	unspanned := make(map[gographt.Vertex]interface{}, vs.Count())
	vs.ForEach(func(vertex gographt.Vertex) error {
		unspanned[vertex] = nil
		return nil
	})

	for len(unspanned) > 0 {
		var root gographt.Vertex
		for vertex, _ := range unspanned {
			root = vertex
			break
		}

		delete(unspanned, root)

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
			if _, found := unspanned[target]; !found {
				target, _ = g.TargetVertexOf(next)
			}
			if _, found := unspanned[target]; !found {
				continue
			}

			mst.es.Add(next)

			delete(unspanned, target)

			edges, _ := g.EdgesOf(target)
			edges.ForEach(func(edge gographt.Edge) error {
				candidate, _ := gographt.OppositeVertexOf(g, edge, target)
				if _, found := unspanned[candidate]; !found {
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
