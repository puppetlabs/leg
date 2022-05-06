# k8sutil

This module provides a framework for declaratively constructing Kubernetes
objects in a controller context as well as various helpful utility functions and
applications.

## controller/app

This package provides a few useful client and controller applications:

* portforward: Programmatic support for port forwarding to a Kubernetes pod or
  service à la `kubectl port-forward`.
* selfsignedsecret: A reconciler loop that generates a self-signed TLS secret.
* tlsproxy: A reverse proxy that wraps an HTTP server with TLS using an
  automatically configured certificate.
* tunnel: A bidirectional stream that connects a locally accessible HTTP server
  to a Kubernetes service inside a cluster.
* webhookcert: A reconciler loop that automatically manages the `caBundle` value
  for a Kubernetes admission webhook.

## controller/eventhandler

This package provides supplementary algorithms to enqueue objects for processing
with a controller that uses controller-runtime.

## controller/obj/lifecycle

This package defines interfaces for entities that conform to a set of standard
behaviors supported by Kubernetes:

* Deletable
* Finalizable
* Labelable/annotatable
* Loadable
* Ownable
* Persistable

## controller/obj/api

This package provides adapters for common builtin Kubernetes object types to the
lifecycle package.

## controller/obj/helper

This package provides helpers that wrap controller-runtime client functionality
to make managing individual Kubernetes objects easier.

## controller/ownerext

This package makes it possible to define ownership guarantees across namespaces.

## manifest

This package supports parsing manifests from YAML sources to Go Kubernetes
objects.

## norm

This package provides normalization routines to make arbitrary string data
conform to common Kubernetes text format requirements, like DNS subdomains.

## test/endtoend

This package wraps the controller-runtime `envtest` package with opinionated
behavior that forces tests to run through an existing cluster. (We recommend using [k3d](https://github.com/rancher/k3d).)

## test/fakeext

This package supplements the Kubernetes client mocks with common reactors.
