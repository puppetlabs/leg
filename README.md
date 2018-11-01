# Insights instrumentation

This package holds instrumentation packages for the services and applications
that run insights.

## Metrics

Insights metrics provides an interface into recording perfomance metrics insight
Go applications and a http.Handler to serve those metrics to 3rd party
monitoring tools such as Prometheus.

Metrics package provides `Must*` functions that log instead of panic. This
allows you to use the library by chaining functions since they return NoOp types
if an error occurs.
