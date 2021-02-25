/*
Package errhandler provides a wrapper reconciler for automatically handling
errors produced by a delegate, which are normally dropped and cause the request
to be requeued.

It also includes errmark-compatible rules for common Kubernetes errors. You can
use these to build up matching actions for the reconciler to handle.
*/
package errhandler
