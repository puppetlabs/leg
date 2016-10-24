package data

type HashMap map[interface{}]interface{}

func (m HashMap) Contains(key interface{}) (found bool) {
	_, found = m[key]
	return
}

func (m HashMap) Put(key, value interface{}) (found bool) {
	found = m.Contains(key)
	m[key] = value

	return
}

func (m HashMap) Get(key interface{}) (value interface{}, found bool) {
	value, found = m[key]
	return
}

func (m HashMap) Remove(key interface{}) (found bool) {
	found = m.Contains(key)
	delete(m, key)

	return
}

func (m HashMap) Empty() bool {
	return m.Size() == 0
}

func (m HashMap) Size() int {
	return len(m)
}

func (m *HashMap) Clear() {
	*m = make(HashMap)
}

func (m *HashMap) Keys() []interface{} {
	return mapKeys(m)
}

func (m *HashMap) Values() []interface{} {
	return mapValues(m)
}

func (m HashMap) ForEach(fn MapIterationFunc) error {
	for key, value := range m {
		if err := fn(key, value); err != nil {
			return err
		}
	}

	return nil
}

func NewHashMap() *HashMap {
	m := make(HashMap)
	return &m
}

func NewHashMapWithCapacity(capacity int) *HashMap {
	m := make(HashMap, capacity)
	return &m
}
