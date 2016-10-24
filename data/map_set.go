package data

type MapBackedSet struct {
	storage Map
}

var mapSetValue = struct{}{}

func (s *MapBackedSet) Contains(elements ...interface{}) bool {
	for _, element := range elements {
		if !s.storage.Contains(element) {
			return false
		}
	}

	return true
}

func (s *MapBackedSet) Add(elements ...interface{}) {
	for _, element := range elements {
		s.storage.Put(element, mapSetValue)
	}
}

func (s *MapBackedSet) Remove(elements ...interface{}) {
	for _, element := range elements {
		s.storage.Remove(element)
	}
}

func (s *MapBackedSet) Empty() bool {
	return s.storage.Empty()
}

func (s *MapBackedSet) Size() int {
	return s.storage.Size()
}

func (s *MapBackedSet) Clear() {
	s.storage.Clear()
}

func (s *MapBackedSet) Values() []interface{} {
	return s.storage.Keys()
}

func (s *MapBackedSet) ForEach(fn SetIterationFunc) error {
	return s.storage.ForEach(func(key, value interface{}) error {
		return fn(key)
	})
}

func NewMapBackedSet(storage Map) *MapBackedSet {
	return &MapBackedSet{storage: storage}
}

func NewHashSet() Set {
	return NewMapBackedSet(NewHashMap())
}

func NewHashSetWithCapacity(capacity int) Set {
	return NewMapBackedSet(NewHashMapWithCapacity(capacity))
}

func NewLinkedHashSet() Set {
	return NewMapBackedSet(NewLinkedHashMap())
}

func NewLinkedHashSetWithCapacity(capacity int) Set {
	return NewMapBackedSet(NewLinkedHashMapWithCapacity(capacity))
}
