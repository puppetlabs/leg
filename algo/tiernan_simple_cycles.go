// Portions of this file are derived from JGraphT, a free Java graph-theory
// library.
//
// (C) Copyright 2013-2017, by Nikolay Ognyanov and Contributors.

package algo

import (
	"fmt"
	"reflect"

	"github.com/reflect/godat"
	"github.com/reflect/gographt"
)

const (
	TiernanSimpleCyclesSupportedFeatures = gographt.DeterministicIteration
)

type TiernanSimpleCycles struct {
	features gographt.GraphFeature
	g        gographt.DirectedGraph
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

func (tsc *TiernanSimpleCycles) Cycles() (cycles [][]gographt.Vertex) {
	if tsc.g.Vertices().Count() == 0 {
		fmt.Printf("out: %v\n", tsc.g.Vertices().Count())
		return
	}

	var path []gographt.Vertex
	pathSet := godat.NewHashSet() // set[]gographt.Vertex

	blocked := make(map[gographt.Vertex]godat.Set) // map[gographt.Vertex]set[]gographt.Vertex
	indices := make(map[gographt.Vertex]int)

	i := 0
	verticesLeft := make([]gographt.Vertex, tsc.g.Vertices().Count())
	tsc.g.Vertices().ForEach(func(vertex gographt.Vertex) error {
		verticesLeft[i] = vertex

		blocked[vertex] = godat.NewHashSet()
		indices[vertex] = i
		i++

		return nil
	})

	var pathStart gographt.Vertex
	pathEnd := verticesLeft[0]
	verticesLeft = verticesLeft[1:]

	path = append(path, pathEnd)
	pathSet.Add(pathEnd)

	for {
		// Path extension.
		for {
			edges, _ := tsc.g.OutgoingEdgesOf(pathEnd)
			err := edges.ForEach(func(edge gographt.Edge) error {
				target, _ := tsc.g.TargetVertexOf(edge)
				if indices[target] > indices[path[0]] && !pathSet.Contains(target) && !blocked[pathEnd].Contains(target) {
					path = append(path, target)
					pathSet.Add(target)
					pathEnd = target
					return godat.ErrStopIteration
				}

				return nil
			})
			if err == godat.ErrStopIteration {
				// We found another extension. Repeat this search.
				continue
			}

			break
		}

		// Circuit confirmation.
		pathStart = path[0]
		if tsc.g.ContainsEdgeBetween(pathEnd, pathStart) {
			cycle := make([]gographt.Vertex, len(path))
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

			path = []gographt.Vertex{pathEnd}
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

func TiernanSimpleCyclesOf(g gographt.DirectedGraph) *TiernanSimpleCycles {
	return &TiernanSimpleCycles{
		features: g.Features() & TiernanSimpleCyclesSupportedFeatures,
		g:        g,
	}
}
