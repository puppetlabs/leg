package main

import (
	"context"
	"log"
	"time"

	"github.com/puppetlabs/insights-instrumentation/metrics"
	"github.com/puppetlabs/insights-instrumentation/metrics/collectors"
	"github.com/puppetlabs/insights-instrumentation/metrics/delegates"
	"github.com/puppetlabs/insights-instrumentation/metrics/server"
)

func main() {
	m, err := metrics.NewNamespace("example", delegates.PrometheusDelegate)
	if err != nil {
		log.Fatal(err)
	}

	copts := collectors.CounterOptions{
		Description: "a count of tasks executed",
		Labels:      []string{"status"},
	}
	if err := m.RegisterCounter("tasks", copts); err != nil {
		log.Fatal(err)
	}

	topts := collectors.TimerOptions{
		Description: "duration of requests to google cloud storage",
		Labels:      []string{"action"},
	}
	if err := m.RegisterTimer("gcs_request_duration", topts); err != nil {
		log.Fatal(err)
	}

	m.OnCounter("tasks", func(c collectors.Counter) {
		c.Add(1)
	})

	c, _ := m.Counter("tasks")

	c.Add(1)
	c.Inc()

	m.OnTimer("gcs_request_duration", func() {
		<-time.After(time.Second)
	})

	m.OnTimer("gcs_request_duration", func() {
		<-time.After(time.Second * 2)
	})

	t, _ := m.Timer("request_timer")

	handle := t.Start()
	<-time.After(time.Millisecond * 500)
	t.ObserveDuration(handle)

	server := server.New(m, server.Options{
		BindAddr: "localhost:2398",
	})

	if err := server.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
