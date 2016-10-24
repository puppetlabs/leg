// Portions of this file are derived from a priority queue implementation
// provided in the Go documentation.
//
// https://golang.org/pkg/container/heap/

package data

import (
	"container/heap"
	"reflect"
)

// An Item is something we manage in a priority queue.
type priorityQueueItem struct {
	value    interface{}
	priority float64
}

// A priorityQueueImpl implements heap.Interface and holds items.
type priorityQueueImpl []*priorityQueueItem

func (pqi priorityQueueImpl) Len() int {
	return len(pqi)
}

func (pqi priorityQueueImpl) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pqi[i].priority > pqi[j].priority
}

func (pqi priorityQueueImpl) Swap(i, j int) {
	pqi[i], pqi[j] = pqi[j], pqi[i]
}

func (pqi *priorityQueueImpl) Push(x interface{}) {
	item := x.(*priorityQueueItem)
	*pqi = append(*pqi, item)
}

func (pqi *priorityQueueImpl) Pop() interface{} {
	old := *pqi
	n := len(old)
	item := old[n-1]
	*pqi = old[0 : n-1]
	return item
}

// We wrap it in a nice interface for the user.

type PriorityQueue struct {
	impl priorityQueueImpl
}

func (pq *PriorityQueue) Add(v interface{}, priority float64) {
	item := &priorityQueueItem{
		value:    v,
		priority: priority,
	}

	heap.Push(&pq.impl, item)
}

func (pq *PriorityQueue) Poll(v interface{}) bool {
	if pq.impl.Len() == 0 {
		return false
	}

	item := heap.Pop(&pq.impl).(*priorityQueueItem)
	reflect.ValueOf(v).Elem().Set(reflect.ValueOf(item.value))

	return true
}

func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{}
	heap.Init(&pq.impl)

	return pq
}
