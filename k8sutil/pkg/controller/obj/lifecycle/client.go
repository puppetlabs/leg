package lifecycle

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type CacheBypasserClient interface {
	client.Client
	BypassCache() client.Reader
}

type cacheBypasserClient struct {
	client.Client
	apiReader client.Reader
}

func (cbc *cacheBypasserClient) BypassCache() client.Reader {
	return cbc.apiReader
}

func CacheBypasserClientForManager(mgr manager.Manager) CacheBypasserClient {
	return &cacheBypasserClient{
		Client:    mgr.GetClient(),
		apiReader: mgr.GetAPIReader(),
	}
}
