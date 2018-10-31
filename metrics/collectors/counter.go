package collectors

type Counter interface {
	Inc()
	Add(float64)
}
