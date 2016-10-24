package data

import (
	"container/list"
)

type linkedHashMapEntry struct {
	key, value interface{}
}

type LinkedHashMap struct {
	accessor *list.List
	storage  map[interface{}]*list.Element
}

func (m *LinkedHashMap) Contains(key interface{}) (found bool) {
	_, found = m.storage[key]
	return
}

func (m *LinkedHashMap) Put(key, value interface{}) (found bool) {
	if _, found = m.storage[key]; !found {
		entry := &linkedHashMapEntry{key: key, value: value}

		e := m.accessor.PushBack(entry)
		m.storage[key] = e
	} else {
		entry := m.storage[key].Value.(*linkedHashMapEntry)
		entry.value = value
	}

	return
}

func (m *LinkedHashMap) Get(key interface{}) (value interface{}, found bool) {
	var e *list.Element

	if e, found = m.storage[key]; found {
		value = e.Value.(*linkedHashMapEntry).value
	}

	return
}

func (m *LinkedHashMap) Remove(key interface{}) (found bool) {
	var e *list.Element

	if e, found = m.storage[key]; found {
		m.accessor.Remove(e)
		delete(m.storage, key)
	}

	return
}

func (m *LinkedHashMap) Empty() bool {
	return m.Size() == 0
}

func (m *LinkedHashMap) Size() int {
	return len(m.storage)
}

func (m *LinkedHashMap) Clear() {
	m.accessor.Init()
	m.storage = make(map[interface{}]*list.Element)
}

func (m *LinkedHashMap) Keys() []interface{} {
	return mapKeys(m)
}

func (m *LinkedHashMap) Values() []interface{} {
	return mapValues(m)
}

func (m *LinkedHashMap) ForEach(fn MapIterationFunc) error {
	for e := m.accessor.Front(); e != nil; e = e.Next() {
		entry := e.Value.(*linkedHashMapEntry)
		if err := fn(entry.key, entry.value); err != nil {
			return err
		}
	}

	return nil
}

func NewLinkedHashMap() *LinkedHashMap {
	return &LinkedHashMap{
		accessor: list.New(),
		storage:  make(map[interface{}]*list.Element),
	}
}

func NewLinkedHashMapWithCapacity(capacity int) *LinkedHashMap {
	return &LinkedHashMap{
		accessor: list.New(),
		storage:  make(map[interface{}]*list.Element, capacity),
	}
}
