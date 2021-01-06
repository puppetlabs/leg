package endtoend

import (
	"context"
	"time"

	"github.com/puppetlabs/leg/lifecycle"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

// EnvironmentOptions are the options that control a given test environment.
type EnvironmentOptions struct {
	CRDDirectoryPaths     []string
	ClientConfigLoader    clientcmd.ClientConfigLoader
	ClientConfigOverrides *clientcmd.ConfigOverrides
	ClientScheme          *runtime.Scheme
}

// EnvironmentOption is a setter for one or more environment options.
type EnvironmentOption interface {
	// ApplyToEnvironmentOptions copies the configuration of this option to the
	// given environment options.
	ApplyToEnvironmentOptions(target *EnvironmentOptions)
}

// ApplyOptions runs each of the given option setters against this environment
// options struct.
func (o *EnvironmentOptions) ApplyOptions(opts []EnvironmentOption) {
	for _, opt := range opts {
		opt.ApplyToEnvironmentOptions(o)
	}
}

// EnvironmentOptionFunc allows a functio to be used as an enviroment option.
type EnvironmentOptionFunc func(target *EnvironmentOptions)

var _ EnvironmentOption = EnvironmentOptionFunc(nil)

// ApplyToEnvironmentOptions copies the configuration of this option to the
// given environment options.
func (eof EnvironmentOptionFunc) ApplyToEnvironmentOptions(target *EnvironmentOptions) {
	eof(target)
}

// EnvironmentWithCRDDirectoryPaths causes any CRDs found in the given
// directories to be loaded into the cluster when the environment connection is
// established.
type EnvironmentWithCRDDirectoryPaths []string

var _ EnvironmentOption = EnvironmentWithCRDDirectoryPaths(nil)

// ApplyToEnvironmentOptions copies the configuration of this option to the
// given environment options.
func (ewcdp EnvironmentWithCRDDirectoryPaths) ApplyToEnvironmentOptions(target *EnvironmentOptions) {
	target.CRDDirectoryPaths = ewcdp
}

// EnvironmentWithClientContext causes the environment connection to use the
// given context from the Kubeconfig(s) instead of the default configured in the
// file.
type EnvironmentWithClientContext string

var _ EnvironmentOption = EnvironmentWithClientContext("")

// ApplyToEnvironmentOptions copies the configuration of this option to the
// given environment options.
func (ewcc EnvironmentWithClientContext) ApplyToEnvironmentOptions(target *EnvironmentOptions) {
	target.ClientConfigOverrides.CurrentContext = string(ewcc)
}

// EnvironmentWithClientKubeconfigs sets the path to one or more Kubeconfig
// files to use to load the connection to the cluster.
//
// This option is incompatible with EnvironmentWithClientConfigLoader.
type EnvironmentWithClientKubeconfigs []string

var _ EnvironmentOption = EnvironmentWithClientKubeconfigs(nil)

// ApplyToEnvironmentOptions copies the configuration of this option to the
// given environment options.
func (ewck EnvironmentWithClientKubeconfigs) ApplyToEnvironmentOptions(target *EnvironmentOptions) {
	lr, ok := target.ClientConfigLoader.(*clientcmd.ClientConfigLoadingRules)
	if !ok {
		lr = &clientcmd.ClientConfigLoadingRules{}
		target.ClientConfigLoader = lr
	}

	lr.Precedence = ewck
	lr.WarnIfAllMissing = true
}

// EnvironmentWithClientConfigLoader replaces the client configuration loader
// with the specified one.
//
// This option is incompatible with EnvironmentWithClientKubeconfigs.
func EnvironmentWithClientConfigLoader(loader clientcmd.ClientConfigLoader) EnvironmentOption {
	return EnvironmentOptionFunc(func(target *EnvironmentOptions) {
		target.ClientConfigLoader = loader
	})
}

// EnvironmentWithClientConfigOverrides replaces the client configuration
// overrides with the specified one.
//
// This option overwrites any changes made by EnvironmentWithClientContext.
func EnvironmentWithClientConfigOverrides(overrides *clientcmd.ConfigOverrides) EnvironmentOption {
	return EnvironmentOptionFunc(func(target *EnvironmentOptions) {
		*target.ClientConfigOverrides = *overrides
	})
}

// EnvironmentWithClientScheme sets the runtime scheme to use for the controller
// and dynamic clients.
func EnvironmentWithClientScheme(s *runtime.Scheme) EnvironmentOption {
	return EnvironmentOptionFunc(func(target *EnvironmentOptions) {
		target.ClientScheme = s
	})
}

// Environment provides an opinionated connection to an existing cluster in
// test. It does not allow a new cluster to be created (like
// controller-runtime's envtest).
type Environment struct {
	Scheme           *runtime.Scheme
	RESTConfig       *rest.Config
	RESTMapper       meta.RESTMapper
	ControllerClient client.Client
	StaticClient     kubernetes.Interface
	DynamicClient    dynamic.Interface

	closer *lifecycle.Closer
}

// Close terminates the connection to this cluster and cleans up resources
// created by the connection.
func (e *Environment) Close(ctx context.Context) error {
	return e.closer.Do(ctx)
}

// NewEnvironment establishes a connection to the cluster given by the relevant
// Kubernetes client configuration and sets up client accessors for it.
func NewEnvironment(opts ...EnvironmentOption) (e *Environment, err error) {
	o := &EnvironmentOptions{
		ClientConfigLoader:    &clientcmd.ClientConfigLoadingRules{},
		ClientConfigOverrides: &clientcmd.ConfigOverrides{},
		ClientScheme:          scheme.Scheme,
	}
	o.ApplyOptions(opts)

	cb := lifecycle.NewCloserBuilder()
	defer func() {
		if err != nil {
			// Best-effort attempt to close anything that needs to be closed.
			_ = cb.Build().Do(context.Background())
		}
	}()

	cfg, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		o.ClientConfigLoader,
		o.ClientConfigOverrides,
	).ClientConfig()
	if err != nil {
		return nil, err
	}

	env := &envtest.Environment{
		Config:                   cfg,
		CRDDirectoryPaths:        o.CRDDirectoryPaths,
		AttachControlPlaneOutput: true,
		UseExistingCluster:       func(b bool) *bool { return &b }(true),
	}

	cfg, err = env.Start()
	if err != nil {
		return nil, err
	}
	cb.Require(env.Stop)

	mapper, err := apiutil.NewDynamicRESTMapper(cfg)
	if err != nil {
		return nil, err
	}

	controllerClient, err := client.New(cfg, client.Options{
		Scheme: o.ClientScheme,
		Mapper: mapper,
	})
	if err != nil {
		return nil, err
	}

	staticClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &Environment{
		Scheme:           o.ClientScheme,
		RESTConfig:       cfg,
		RESTMapper:       mapper,
		ControllerClient: controllerClient,
		StaticClient:     staticClient,
		DynamicClient:    dynamicClient,

		closer: cb.Build(),
	}, nil
}

// WithEnvironment runs a function with a connection to a cluster and then tears
// down the connection.
func WithEnvironment(opts []EnvironmentOption, fn func(e *Environment)) (err error) {
	e, err := NewEnvironment(opts...)
	if err != nil {
		return err
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err = e.Close(ctx)
	}()

	fn(e)
	return
}
