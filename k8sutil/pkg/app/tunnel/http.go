package tunnel

import (
	"context"
	"net/url"

	corev1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/corev1"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HTTPImage is the Inlets Docker image to use for tunneling.
const HTTPImage = "ghcr.io/puppetlabs/inlets:latest"

// HTTP represents the resources required to maintain a tunnel.
type HTTP struct {
	Key client.ObjectKey

	OwnerConfigMap *corev1obj.ConfigMap
	Service        *corev1obj.Service
	Pod            *corev1obj.Pod
}

var _ lifecycle.Deleter = &HTTP{}
var _ lifecycle.Loader = &HTTP{}
var _ lifecycle.Ownable = &HTTP{}
var _ lifecycle.Persister = &HTTP{}

func (h *HTTP) Delete(ctx context.Context, cl client.Client, opts ...lifecycle.DeleteOption) (bool, error) {
	return h.OwnerConfigMap.Delete(ctx, cl, opts...)
}

func (h *HTTP) Load(ctx context.Context, cl client.Client) (bool, error) {
	return lifecycle.Loaders{
		h.OwnerConfigMap,
		h.Service,
		h.Pod,
	}.Load(ctx, cl)
}

func (h *HTTP) Owned(ctx context.Context, owner lifecycle.TypedObject) error {
	return h.OwnerConfigMap.Owned(ctx, owner)
}

func (h *HTTP) Persist(ctx context.Context, cl client.Client) error {
	return lifecycle.OwnershipPersister{
		Owner: h.OwnerConfigMap,
		Dependent: lifecycle.OwnablePersisters{
			h.Service,
			h.Pod,
		},
	}.Persist(ctx, cl)
}

// URL returns the HTTP URL to the tunnel service.
func (h *HTTP) URL() string {
	return (&url.URL{
		Scheme: "http",
		Host:   h.Service.DNSName(),
	}).String()
}

// NewHTTP creates a new, unconfigured HTTP-only tunnel service.
func NewHTTP(key client.ObjectKey) *HTTP {
	return &HTTP{
		Key: key,

		OwnerConfigMap: corev1obj.NewConfigMap(key),
		Service:        corev1obj.NewService(key),
		Pod:            corev1obj.NewPod(key),
	}
}

// ConfigureHTTP sets up the tunnel to receive connections. Use
// WithHTTPConnection to connect to the tunnel after persisting a configured
// tunnel.
func ConfigureHTTP(h *HTTP) *HTTP {
	selector := map[string]string{
		"app.kubernetes.io/name":     "tunnel.http",
		"app.kubernetes.io/instance": h.Key.Name,
	}

	h.Service.Object.Spec = corev1.ServiceSpec{
		Type: corev1.ServiceTypeClusterIP,
		Ports: []corev1.ServicePort{
			{
				Name:       "proxy-http",
				TargetPort: intstr.FromString("proxy-http"),
				Protocol:   corev1.ProtocolTCP,
				Port:       80,
			},
		},
		Selector: selector,
	}

	h.Pod.Object.ObjectMeta.Labels = selector
	h.Pod.Object.Spec = corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:  "tunnel",
				Image: HTTPImage,
				Args: []string{
					"server",
					"--port", "8000",
					"--control-port", "8080",
				},
				Ports: []corev1.ContainerPort{
					{
						Name:          "proxy-http",
						ContainerPort: 8000,
						Protocol:      corev1.ProtocolTCP,
					},
					{
						Name:          "tunnel",
						ContainerPort: 8080,
						Protocol:      corev1.ProtocolTCP,
					},
				},
				LivenessProbe: &corev1.Probe{
					Handler: corev1.Handler{
						TCPSocket: &corev1.TCPSocketAction{
							Port: intstr.FromString("tunnel"),
						},
					},
				},
			},
		},
	}

	return h
}

// ApplyHTTP loads, configures, and persists any changes to an HTTP tunnel using
// the given client.
func ApplyHTTP(ctx context.Context, cl client.Client, key client.ObjectKey) (*HTTP, error) {
	h := NewHTTP(key)

	if _, err := h.Load(ctx, cl); err != nil {
		return nil, err
	}

	h = ConfigureHTTP(h)

	if err := h.Persist(ctx, cl); err != nil {
		return nil, err
	}

	return h, nil
}
