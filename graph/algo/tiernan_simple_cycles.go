// Portions of this file are derived from JGraphT, a free Java graph-theory
// library.
//
// (C) Copyright 2013-2017, by Nikolay Ognyanov and Contributors.

package algo

import (
	"reflect"

	"github.com/puppetlabs/horsehead/v2/datastructure"
	"github.com/puppetlabs/horsehead/v2/graph"
)

const (
	TiernanSimpleCyclesSupportedFeatures = graph.DeterministicIteration
)

type TiernanSimpleCycles struct {
	features graph.GraphFeature
	g        graph.DirectedGraph
}

func (tsc *TiernanSimpleCycles) CyclesInto(into interface{}) {
	p := reflect.ValueOf(into).Elem()
	slice := p

	for _, cycle := range tsc.Cycles() {
		size := len(cycle)
		is := reflect.MakeSlice(p.Type().Elem(), size, size)
		ist := is.Type().Elem()

		for i, vertex := range cycle {
			v := reflect.ValueOf(vertex)
			if !v.IsValid() {
				v = reflect.Zero(ist)
			}

			is.Index(i).Set(v)
		}

		slice = reflect.Append(slice, is)
	}

	p.Set(slice)
}

func (tsc *TiernanSimpleCycles) Cycles() (cycles [][]graph.Vertex) {
	if tsc.g.Vertices().Count() == 0 {
		return
	}

	var path []graph.Vertex
	pathSet := datastructure.NewHashSet() // set[]graph.Vertex

	blocked := make(map[graph.Vertex]datastructure.Set) // map[graph.Vertex]set[]graph.Vertex
	indices := make(map[graph.Vertex]int)

	i := 0
	verticesLeft := make([]graph.Vertex, tsc.g.Vertices().Count())
	tsc.g.Vertices().ForEach(func(vertex graph.Vertex) error {
		verticesLeft[i] = vertex

		blocked[vertex] = datastructure.NewHashSet()
		indices[vertex] = i
		i++

		return nil
	})

	var pathStart graph.Vertex
	pathEnd := verticesLeft[0]
	verticesLeft = verticesLeft[1:]

	path = append(path, pathEnd)
	pathSet.Add(pathEnd)

	for {
		// Path extension.
		for {
			edges, _ := tsc.g.OutgoingEdgesOf(pathEnd)
			err := edges.ForEach(func(edge graph.Edge) error {
				target, _ := tsc.g.TargetVertexOf(edge)
				if indices[target] > indices[path[0]] && !pathSet.Contains(target) && !blocked[pathEnd].Contains(target) {
					path = append(path, target)
					pathSet.Add(target)
					pathEnd = target
					return datastructure.ErrStopIteration
				}

				return nil
			})
			if err == datastructure.ErrStopIteration {
				// We found another extension. Repeat this search.
				continue
			}

			break
		}

		// Circuit confirmation.
		pathStart = path[0]
		if tsc.g.ContainsEdgeBetween(pathEnd, pathStart) {
			cycle := make([]graph.Vertex, len(path))
			copy(cycle, path)

			cycles = append(cycles, cycle)
		}

		// Vertex closure.
		if len(path) > 1 {
			path = path[:len(path)-1]
			pathSet.Remove(pathEnd)
			blocked[pathEnd].Clear()

			pathLast := pathEnd
			pathEnd = path[len(path)-1]
			blocked[pathEnd].Add(pathLast)

			continue
		}

		// Advance initial index.
		if len(verticesLeft) > 0 {
			pathEnd = verticesLeft[0]
			verticesLeft = verticesLeft[1:]

			path = []graph.Vertex{pathEnd}
			pathSet.Clear()
			pathSet.Add(pathEnd)

			for _, blockedPath := range blocked {
				blockedPath.Clear()
			}

			continue
		}

		// Terminate.
		break
	}

	return
}

func TiernanSimpleCyclesOf(g graph.DirectedGraph) *TiernanSimpleCycles {
	return &TiernanSimpleCycles{
		features: g.Features() & TiernanSimpleCyclesSupportedFeatures,
		g:        g,
	}
}
