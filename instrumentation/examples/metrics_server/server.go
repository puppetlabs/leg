package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/puppetlabs/horsehead/v2/instrumentation/metrics"
	"github.com/puppetlabs/horsehead/v2/instrumentation/metrics/collectors"
	"github.com/puppetlabs/horsehead/v2/instrumentation/metrics/delegates"
	"github.com/puppetlabs/horsehead/v2/instrumentation/metrics/server"
)

const (
	namespace           = "example"
	storageBackendCall  = "storage_backend_call"
	taskCount           = "task_count"
	requestCount        = "request_count"
	httpRequestDuration = "http_handler_duration"
)

type testApp struct {
	m *metrics.Metrics
}

func (t testApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Calls MustCounter to get a metric registered as requestCounter (in this case it's name is request_count)
	// and applys the label "method" with the request method.
	c := t.m.MustCounter(requestCount, metrics.NewLabel("method", r.Method))
	// Here we are using Add(1), but c.Inc() would also increment the count by one.
	c.Add(1)

	// Now we are getting a metric registered as storageBackendCall to track the duration of a function call
	// that does some work with a storage backend. We are setting the labels backend=gcs and action=put to give
	// the metric depth, allowing us to layer charts with more context about operations.
	gcsTimer := t.m.MustTimer(storageBackendCall,
		metrics.NewLabel("backend", "gcs"),
		metrics.NewLabel("action", "put"),
	)

	// Using the timer we got from MustTimer, we are passing it into the OnTimer function that takes a callback.
	// This wraps the func() to get a duration of how long it took to run by calling Start() on the timer, then calling
	// func(), then calling ObserveDuration().
	t.m.OnTimer(gcsTimer, func() {
		random := rand.Intn(10)
		<-time.After(time.Millisecond * time.Duration(random))
	})

	// This is during the same thing as the gcsTimer, but it's using the labels backend=s3 and action=get instead.
	s3Timer := t.m.MustTimer(storageBackendCall,
		metrics.NewLabel("backend", "s3"),
		metrics.NewLabel("action", "get"),
	)

	t.m.OnTimer(s3Timer, func() {
		random := rand.Intn(10)
		<-time.After(time.Millisecond * time.Duration(random))
	})

	// This is demoing how we can use the gcsTimer metric without OnTimer and also how we can change the labels
	// once we want to ObserveDuration()
	h := gcsTimer.Start()
	random := rand.Intn(10)
	<-time.After(time.Millisecond * time.Duration(random))
	gcsTimer.ObserveDuration(h, metrics.NewLabel("backend", "backblaze"), metrics.NewLabel("action", "post"))
}

func main() {
	// We start bootstrapping our metrics collection by creating a new metrics.Metric object
	// that represents a namespace. This is essencially bucketizing this specific system in the
	// metrics storage backend. The namespace must only include characters a-z and _.
	m, err := metrics.NewNamespace(namespace, metrics.Options{
		DelegateType:  delegates.PrometheusDelegate,
		ErrorBehavior: metrics.ErrorBehaviorLog,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Now we are registering a Counter metric as requestCount. All metrics that you want to use need to be
	// registered with the metrics.Metric object.
	//
	// Here we are calling RegisterCounter, which could return an error. This error must be handled by the caller.
	copts := collectors.CounterOptions{
		Description: "a count of requests executed",
		Labels:      []string{"method"},
	}
	if err := m.RegisterCounter(requestCount, copts); err != nil {
		log.Fatal(err)
	}

	// Now we are registering another Counter, but using MustRegisterCounter. This will catch any error returned
	// and pass it to the internal metrics.handleError method, which can be configured to handle the error automatically.
	// The default error handling behavior is to log the error and return a noop metric, which will do nothing.
	m.MustRegisterCounter(taskCount, collectors.CounterOptions{
		Description: "a count of tasks executed",
		Labels:      []string{"status"},
	})

	topts := collectors.TimerOptions{
		Description: "duration of requests to google cloud storage",
		Labels:      []string{"backend", "action"},
	}
	if err := m.RegisterTimer(storageBackendCall, topts); err != nil {
		log.Fatal(err)
	}

	hopts := collectors.DurationMiddlewareOptions{
		Description: "http request duration middleware",
		Labels:      []string{"handler", "code", "method"},
	}
	if err := m.RegisterDurationMiddleware(httpRequestDuration, hopts); err != nil {
		log.Fatal(err)
	}

	// This is an example of pre-configuring metrics with label values. This is useful if your label values
	// are constant (which they should be due to heavy resource usage with high cardinality data). You can create
	// a map if these labeled metrics and grab the one you want to observe with a key from a list of consts.
	failedTasks := m.MustCounter(taskCount, metrics.NewLabel("status", "failed"))
	failedTasks.Add(5)

	succeededTasks := m.MustCounter(taskCount, metrics.NewLabel("status", "succeeded"))
	succeededTasks.Add(1)
	succeededTasks.Inc()

	// This one creates a labeled metric that is a HTTP middleware handler that will track the duration
	// of all HTTP handlers wrapped by it. If the URL contains user IDs or something, don't add a label
	// that has the path in it. This will cause most metric backends to explode.
	mw := m.MustDurationMiddleware(httpRequestDuration, metrics.NewLabel("handler", "test_handler"))
	mux := http.NewServeMux()
	mux.Handle("/", mw.Wrap(testApp{m}))

	log.Println("example server running at: http://localhost:2399/")
	go http.ListenAndServe("localhost:2399", mux)

	// Now that we have registered a bunch of metrics and configured them, we now create a metrics.Server
	// (which is just a HTTP server) and tell it to serve the metrics on http://locahost:3050/
	server := server.New(m, server.Options{
		BindAddr: "localhost:3050",
		Path:     "/",
	})

	log.Println("metrics server running at: http://localhost:3050/")
	if err := server.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
