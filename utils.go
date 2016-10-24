package gographt

import (
	"sync/atomic"
)

var defaultEdgeIndex uint64

type defaultEdge struct {
	index uint64
}

func NewEdge() Edge {
	return &defaultEdge{atomic.AddUint64(&defaultEdgeIndex, 1)}
}

func OppositeVertexOf(g Graph, e Edge, v Vertex) (Vertex, error) {
	test, err := g.SourceVertexOf(e)
	if err != nil {
		return nil, err
	}

	if test != v {
		return test, nil
	} else {
		return g.TargetVertexOf(e)
	}
}
