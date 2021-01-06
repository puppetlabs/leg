// Package tlsproxy provides a reverse proxy that wraps an HTTP service with
// HTTPS. It uses Square's Ghostunnel (https://github.com/ghostunnel/ghostunnel)
// for the proxy.
//
// Certificate generation is currently hardcoded; custom certificates cannot be
// supplied. Certificates are valid for a short time making this package mostly
// suitable for testing. A certificate authority is generated and provided.
package tlsproxy

import (
	"context"
	"encoding/pem"
	"net/url"

	appsv1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/appsv1"
	corev1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/corev1"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Image is the Docker image for Ghostunnel to use.
const Image = "squareup/ghostunnel:v1.5.2"

// TLSProxy is the service and deployment information for an instance of a
// reverse proxy.
type TLSProxy struct {
	Key client.ObjectKey

	OwnerConfigMap *corev1obj.ConfigMap
	Service        *corev1obj.Service
	Secret         *corev1obj.TLSSecret
	Deployment     *appsv1obj.Deployment
}

var _ lifecycle.Deleter = &TLSProxy{}
var _ lifecycle.Loader = &TLSProxy{}
var _ lifecycle.Ownable = &TLSProxy{}
var _ lifecycle.Persister = &TLSProxy{}

func (tp *TLSProxy) Delete(ctx context.Context, cl client.Client, opts ...lifecycle.DeleteOption) (bool, error) {
	return tp.OwnerConfigMap.Delete(ctx, cl, opts...)
}

func (tp *TLSProxy) Load(ctx context.Context, cl client.Client) (bool, error) {
	return lifecycle.Loaders{
		tp.OwnerConfigMap,
		tp.Service,
		tp.Secret,
		tp.Deployment,
	}.Load(ctx, cl)
}

func (tp *TLSProxy) Owned(ctx context.Context, owner lifecycle.TypedObject) error {
	return tp.OwnerConfigMap.Owned(ctx, owner)
}

func (tp *TLSProxy) Persist(ctx context.Context, cl client.Client) error {
	return lifecycle.OwnershipPersister{
		Owner: tp.OwnerConfigMap,
		Dependent: lifecycle.OwnablePersisters{
			tp.Service,
			tp.Secret,
			tp.Deployment,
		},
	}.Persist(ctx, cl)
}

// URL returns the HTTPS URL to the wrapped service.
func (tp *TLSProxy) URL() string {
	return (&url.URL{
		Scheme: "https",
		Host:   tp.Service.DNSName(),
	}).String()
}

// CertificateAuthorityPEM returns the PEM-encoded CA to use for validating
// connections to the service.
func (tp *TLSProxy) CertificateAuthorityPEM() ([]byte, error) {
	cert, err := tp.Secret.Certificate()
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Certificate[len(cert.Certificate)-1],
	}), nil
}

// New creates a new, unconfigured TLS proxy with the given namespace and name.
func New(key client.ObjectKey) *TLSProxy {
	return &TLSProxy{
		Key: key,

		OwnerConfigMap: corev1obj.NewConfigMap(key),
		Service:        corev1obj.NewService(key),
		Secret:         corev1obj.NewTLSSecret(key),
		Deployment:     appsv1obj.NewDeployment(key),
	}
}

// Configure sets up the TLS proxy to connect to the given upstream address.
func Configure(tp *TLSProxy, addr string) (*TLSProxy, error) {
	selector := map[string]string{
		"app.kubernetes.io/name":     "tls-proxy",
		"app.kubernetes.io/instance": tp.Key.Name,
	}

	tp.Service.Object.Spec = corev1.ServiceSpec{
		Type: corev1.ServiceTypeClusterIP,
		Ports: []corev1.ServicePort{
			{
				Name:       "proxy-https",
				TargetPort: intstr.FromString("proxy-https"),
				Protocol:   corev1.ProtocolTCP,
				Port:       443,
			},
		},
		Selector: selector,
	}

	if tp.Secret.Object.Data["tls.key"] == nil || tp.Secret.Object.Data["tls.crt"] == nil {
		bundle, err := generateTLSCertificateBundle(tp.Service.DNSName())
		if err != nil {
			return nil, err
		}

		tp.Secret.Object.Data = map[string][]byte{
			"tls.key": bundle.ServerKeyPEM,
			"tls.crt": bundle.BundlePEM,
			"ca.crt":  bundle.AuthorityPEM,
		}
	}

	tp.Deployment.Object.Spec = appsv1.DeploymentSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: selector,
		},
		Strategy: appsv1.DeploymentStrategy{
			Type: appsv1.RollingUpdateDeploymentStrategyType,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: selector,
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "proxy",
						Image: Image,
						Args: []string{
							"server",
							"--listen", ":9000",
							"--unsafe-target",
							"--target", addr,
							"--disable-authentication",
							"--key", "/var/run/tls-proxy/tls/tls.key",
							"--cert", "/var/run/tls-proxy/tls/tls.crt",
							"--status", ":9001",
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "tls",
								MountPath: "/var/run/tls-proxy/tls",
							},
						},
						Ports: []corev1.ContainerPort{
							{
								Name:          "proxy-https",
								ContainerPort: 9000,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								Name:          "proxy-status",
								ContainerPort: 9001,
								Protocol:      corev1.ProtocolTCP,
							},
						},
						ReadinessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path:   "/_status",
									Port:   intstr.FromString("proxy-status"),
									Scheme: corev1.URISchemeHTTPS,
								},
							},
						},
					},
				},
				Volumes: []corev1.Volume{
					{
						Name: "tls",
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: tp.Secret.Key.Name,
							},
						},
					},
				},
			},
		},
	}

	return tp, nil
}

// Apply loads, configures, and persists any changes to a TLS proxy using the
// given client.
func Apply(ctx context.Context, cl client.Client, key client.ObjectKey, addr string) (*TLSProxy, error) {
	tp := New(key)

	if _, err := tp.Load(ctx, cl); err != nil {
		return nil, err
	}

	tp, err := Configure(tp, addr)
	if err != nil {
		return nil, err
	}

	if err := tp.Persist(ctx, cl); err != nil {
		return nil, err
	}

	return tp, nil
}
