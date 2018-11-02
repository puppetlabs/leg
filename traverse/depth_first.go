// Portions of this file are derived from JGraphT, a free Java graph-theory
// library.
//
// (C) Copyright 2003-2018, by Liviu Rau and Contributors.

package traverse

import (
	"reflect"

	"github.com/reflect/gographt"
)

var depthFirstSentinel = struct{}{}

type depthFirstStackElement struct {
	vertex gographt.Vertex

	prevStackElement *depthFirstStackElement
}

type DepthFirstTraverser struct {
	g     gographt.Graph
	start gographt.Vertex
}

func (t *DepthFirstTraverser) forEachEdgeOf(vertex gographt.Vertex, fn gographt.EdgeSetIterationFunc) {
	var edges gographt.EdgeSet
	if dg, ok := t.g.(gographt.DirectedGraph); ok {
		edges, _ = dg.OutgoingEdgesOf(vertex)
	} else {
		edges, _ = t.g.EdgesOf(vertex)
	}

	edges.ForEach(fn)
}

func (t *DepthFirstTraverser) ForEach(fn func(vertex gographt.Vertex) error) error {
	seen := make(map[gographt.Vertex]struct{})

	stack := &depthFirstStackElement{vertex: t.start}
	for stack != nil {
		var cur gographt.Vertex
		cur, stack = stack.vertex, stack.prevStackElement

		if _, found := seen[cur]; found {
			continue
		}

		seen[cur] = depthFirstSentinel

		if err := fn(cur); err != nil {
			return err
		}

		t.forEachEdgeOf(cur, func(edge gographt.Edge) error {
			next, _ := gographt.OppositeVertexOf(t.g, edge, cur)

			stack = &depthFirstStackElement{
				vertex:           next,
				prevStackElement: stack,
			}

			return nil
		})
	}

	return nil
}

func (t *DepthFirstTraverser) ForEachInto(fn interface{}) error {
	fnr := reflect.ValueOf(fn)
	fnt := fnr.Type()

	if fnt.NumOut() != 1 {
		panic(ErrInvalidFuncSignature)
	}

	return t.ForEach(func(vertex gographt.Vertex) error {
		p := reflect.ValueOf(vertex)
		if !p.IsValid() {
			p = reflect.Zero(fnt.In(0))
		}

		r := fnr.Call([]reflect.Value{p})

		err := r[0]
		if !err.IsNil() {
			return err.Interface().(error)
		}

		return nil
	})
}

func NewDepthFirstTraverser(g gographt.Graph, start gographt.Vertex) *DepthFirstTraverser {
	return &DepthFirstTraverser{
		g:     g,
		start: start,
	}
}
