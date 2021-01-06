package fakeext

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/testing"
)

// Fake is an interface to wrap client-go/testing.Fake, which is embedded into
// many fake clientsets.
type Fake interface {
	Actions() []testing.Action
	AddProxyReactor(resource string, reaction testing.ProxyReactionFunc)
	AddReactor(verb, resource string, reaction testing.ReactionFunc)
	AddWatchReactor(resource string, reaction testing.WatchReactionFunc)
	ClearActions()
	Invokes(action testing.Action, defaultReturnObj runtime.Object) (runtime.Object, error)
	InvokesProxy(action testing.Action) rest.ResponseWrapper
	InvokesWatch(action testing.Action) (watch.Interface, error)
	PrependProxyReactor(resource string, reaction testing.ProxyReactionFunc)
	PrependReactor(verb, resource string, reaction testing.ReactionFunc)
	PrependWatchReactor(resource string, reaction testing.WatchReactionFunc)
}

// Extension allows for predefined customization of a mock Kubernetes client.
type Extension interface {
	// OnInitial runs over the set of initial objects provided to the mock
	// client.
	OnInitial(objs []runtime.Object) error

	// OnNewFake runs when the mock client is actually created.
	OnNewFake(f Fake) error
}

// NewKubernetesClientsetWithExtensions creates a client-go/kubernetes.Interface
// mock with the given extensions.
func NewKubernetesClientsetWithExtensions(objs []runtime.Object, exts []Extension) (kubernetes.Interface, error) {
	for _, ext := range exts {
		if err := ext.OnInitial(objs); err != nil {
			return nil, err
		}
	}

	kc := fake.NewSimpleClientset(objs...)
	for _, ext := range exts {
		if err := ext.OnNewFake(kc); err != nil {
			return nil, err
		}
	}

	return kc, nil
}

// NewDynamicClientWithExtensions creates a client-go/dynamic.Interface mock
// with the given extensions.
func NewDynamicClientWithExtensions(scheme *runtime.Scheme, objs []runtime.Object, exts []Extension) (dynamic.Interface, error) {
	for _, ext := range exts {
		if err := ext.OnInitial(objs); err != nil {
			return nil, err
		}
	}

	kc := dynamicfake.NewSimpleDynamicClient(scheme, objs...)
	for _, ext := range exts {
		if err := ext.OnNewFake(kc); err != nil {
			return nil, err
		}
	}

	return kc, nil
}
