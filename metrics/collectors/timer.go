package collectors

type TimerHandle struct{}

type TimerOptions struct {
	HistogramBoundaries []float64
}

type Timer interface {
	Start() *TimerHandle
	ObserveDuration(*TimerHandle)
}
