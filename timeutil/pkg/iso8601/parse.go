package iso8601

type parseFunc func(string) (string, error)

func parse(s string, fns []parseFunc) (rest string, err error) {
	rest = s

	for _, fn := range fns {
		rest, err = fn(rest)
		if err != nil {
			return
		}
	}

	return
}
