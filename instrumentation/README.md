# Instrumentation

This package holds instrumentation packages (Sentry for alerting and Prometheus for metrics).

## Metrics

Insights metrics provides an interface into recording perfomance metrics insight
Go applications and a http.Handler to serve those metrics to 3rd party
monitoring tools such as Prometheus.

Metrics package provides `Must*` functions that log instead of panic. This
allows you to use the library by chaining functions since they return NoOp types
if an error occurs.
