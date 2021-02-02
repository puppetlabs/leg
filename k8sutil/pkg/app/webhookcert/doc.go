/*
Package webhookcert automatically manages the caBundle value for a Kubernetes
admission webhook.

It is inspired by Knative's webhook package
(https://github.com/knative/pkg/tree/master/webhook) but does not depend on the
rest of the Knative ecosystem, instead integrating with controller-runtime.

This package exposes a method, AddReconcilerToManager, that should be used with
an already-instantiated controller-runtime Manager to add this automation to
your existing controller.
*/
package webhookcert
