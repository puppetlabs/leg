package gographt

import (
	"github.com/reflect/gographt/data"
)

type UnenforcedSliceEdgeSet []Edge

func (es UnenforcedSliceEdgeSet) Contains(edge Edge) bool {
	for _, e := range es {
		if e == edge {
			return true
		}
	}

	return false
}

func (es UnenforcedSliceEdgeSet) Count() uint {
	return uint(len(es))
}

func (es UnenforcedSliceEdgeSet) AsSlice() []Edge {
	return es
}

func (es UnenforcedSliceEdgeSet) ForEach(fn EdgeSetIterationFunc) error {
	for _, edge := range es {
		if err := fn(edge); err != nil {
			return err
		}
	}

	return nil
}

func (es *UnenforcedSliceEdgeSet) Add(edge Edge) {
	*es = append(*es, edge)
}

type edgeSet struct {
	storage data.Set
}

func (es *edgeSet) Contains(edge Edge) bool {
	return es.storage.Contains(edge)
}

func (es *edgeSet) Count() uint {
	return uint(es.storage.Size())
}

func (es *edgeSet) AsSlice() []Edge {
	s := make([]Edge, es.Count())

	i := 0
	es.ForEach(func(edge Edge) error {
		s[i] = edge
		i++

		return nil
	})

	return s
}

func (es *edgeSet) ForEach(fn EdgeSetIterationFunc) error {
	return es.storage.ForEach(func(element interface{}) error {
		return fn(element.(Edge))
	})
}

func (es *edgeSet) Add(edge Edge) {
	es.storage.Add(edge)
}

func (es *edgeSet) Remove(edge Edge) {
	es.storage.Remove(edge)
}

func NewMutableEdgeSet(storage data.Set) MutableEdgeSet {
	return &edgeSet{storage}
}
