package data

func mapKeys(m Map) []interface{} {
	keys := make([]interface{}, m.Size())

	i := 0
	m.ForEach(func(key, value interface{}) error {
		keys[i] = key
		i++

		return nil
	})

	return keys
}

func mapValues(m Map) []interface{} {
	values := make([]interface{}, m.Size())

	i := 0
	m.ForEach(func(key, value interface{}) error {
		values[i] = value
		i++

		return nil
	})

	return values
}
