package tunnel

import (
	"context"
	"net"
	"strconv"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/app/tlsproxy"
	corev1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/corev1"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HTTPS is a helper that manages the combination of an HTTP tunnel with a TLS
// proxy.
type HTTPS struct {
	Key client.ObjectKey

	OwnerConfigMap *corev1obj.ConfigMap
	HTTP           *HTTP
	TLSProxy       *tlsproxy.TLSProxy
}

var _ lifecycle.Deleter = &HTTPS{}
var _ lifecycle.Loader = &HTTPS{}
var _ lifecycle.Ownable = &HTTPS{}
var _ lifecycle.Persister = &HTTPS{}

func (h *HTTPS) Delete(ctx context.Context, cl client.Client) (bool, error) {
	return h.OwnerConfigMap.Delete(ctx, cl)
}

func (h *HTTPS) Load(ctx context.Context, cl client.Client) (bool, error) {
	return lifecycle.Loaders{
		h.OwnerConfigMap,
		h.HTTP,
		h.TLSProxy,
	}.Load(ctx, cl)
}

func (h *HTTPS) Owned(ctx context.Context, owner lifecycle.TypedObject) error {
	return h.OwnerConfigMap.Owned(ctx, owner)
}

func (h *HTTPS) Persist(ctx context.Context, cl client.Client) error {
	return lifecycle.OwnershipPersister{
		Owner: h.OwnerConfigMap,
		Dependent: lifecycle.OwnablePersisters{
			h.HTTP,
			h.TLSProxy,
		},
	}.Persist(ctx, cl)
}

// URL returns the HTTPS URL to the wrapped tunnel.
func (h *HTTPS) URL() string {
	return h.TLSProxy.URL()
}

// CertificateAuthorityPEM returns the PEM-encoded CA to use for validating
// connections to the tunnel from within the cluster.
func (h *HTTPS) CertificateAuthorityPEM() ([]byte, error) {
	return h.TLSProxy.CertificateAuthorityPEM()
}

// NewHTTPS creates an unconfigured tunnel and TLS proxy pair.
func NewHTTPS(key client.ObjectKey) *HTTPS {
	return &HTTPS{
		Key: key,

		OwnerConfigMap: corev1obj.NewConfigMap(key),
		HTTP:           NewHTTP(helper.SuffixObjectKeyName(key, "tunnel")),
		TLSProxy:       tlsproxy.New(helper.SuffixObjectKeyName(key, "proxy")),
	}
}

// ConfigureHTTPS sets up the tunnel and TLS proxy.
func ConfigureHTTPS(h *HTTPS) (*HTTPS, error) {
	h.HTTP = ConfigureHTTP(h.HTTP)

	addr := net.JoinHostPort(
		h.HTTP.Service.DNSName(),
		strconv.FormatInt(int64(h.HTTP.Service.Object.Spec.Ports[0].Port), 10),
	)

	var err error
	h.TLSProxy, err = tlsproxy.Configure(h.TLSProxy, addr)
	if err != nil {
		return nil, err
	}

	return h, nil
}

// ApplyHTTPS loads, configures, and persists any changes to an HTTP tunnel with
// TLS reverse proxy using the given client.
func ApplyHTTPS(ctx context.Context, cl client.Client, key client.ObjectKey) (*HTTPS, error) {
	h := NewHTTPS(key)

	if _, err := h.Load(ctx, cl); err != nil {
		return nil, err
	}

	h, err := ConfigureHTTPS(h)
	if err != nil {
		return nil, err
	}

	if err := h.Persist(ctx, cl); err != nil {
		return nil, err
	}

	return h, nil
}
