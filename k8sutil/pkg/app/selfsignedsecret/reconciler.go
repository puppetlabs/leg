package selfsignedsecret

import (
	"context"
	"crypto/x509"
	"fmt"
	"time"

	corev1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/corev1"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	"github.com/puppetlabs/leg/k8sutil/pkg/internal/tls"
	"github.com/puppetlabs/leg/stringutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;update

type Reconciler struct {
	cl           client.Client
	organization string
	dnsNames     []string
}

func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (res reconcile.Result, err error) {
	klog.InfoS("self-signed TLS secret reconciler: starting reconcile", "secret", req.NamespacedName)
	defer klog.InfoS("self-signed TLS secret reconciler: ending reconcile", "secret", req.NamespacedName)
	defer func() {
		if err != nil {
			klog.ErrorS(err, "self-signed TLS secret reconciler: failed to reconcile", "secret", req.NamespacedName)
		}
	}()

	secret := corev1obj.NewTLSSecret(req.NamespacedName)
	if _, err := (lifecycle.RequiredLoader{Loader: secret}).Load(ctx, r.cl); err != nil {
		return reconcile.Result{}, err
	}

	if secretCertificateValid(secret, r.dnsNames) {
		return reconcile.Result{RequeueAfter: 1 * time.Hour}, nil
	}

	bundle, err := tls.GenerateCertificateBundle(
		tls.CertificateBundleWithOrganization(r.organization),
		tls.CertificateBundleWithCACommonName("Kubernetes Self-Signed TLS Secret CA"),
		tls.CertificateBundleWithCertificateDNSNames(r.dnsNames),
		tls.CertificateBundleWithValidity(90*24*time.Hour),
	)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to generate certificate: %w", err)
	}

	secret.Object.Data = map[string][]byte{
		"tls.key": bundle.ServerKeyPEM,
		"tls.crt": bundle.BundlePEM,
		"ca.crt":  bundle.AuthorityPEM,
	}

	if err := secret.Persist(ctx, r.cl); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to persist Secret %q: %w", secret.Key, err)
	}

	return reconcile.Result{RequeueAfter: 1 * time.Hour}, nil
}

func secretCertificateValid(secret *corev1obj.TLSSecret, expectedDNSNames []string) bool {
	cert, err := secret.Certificate()
	if err != nil {
		return false
	}

	if len(cert.Certificate) < 2 {
		return false
	}

	data, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return false
	}

	if l, r := stringutil.Diff(data.DNSNames, expectedDNSNames); len(l) != 0 || len(r) != 0 {
		return false
	}

	if time.Now().Add(24 * time.Hour).After(data.NotAfter) {
		return false
	}

	return true
}

func NewReconciler(cl client.Client, organization string, dnsNames ...string) *Reconciler {
	return &Reconciler{
		cl:           cl,
		organization: organization,
		dnsNames:     dnsNames,
	}
}

func AddReconcilerToManager(mgr manager.Manager, secretKey client.ObjectKey, organization string, dnsNames ...string) error {
	r := NewReconciler(mgr.GetClient(), organization, dnsNames...)

	return builder.ControllerManagedBy(mgr).
		For(&corev1.Secret{}, builder.WithPredicates(predicate.NewPredicateFuncs(func(obj client.Object) bool {
			return client.ObjectKeyFromObject(obj) == secretKey
		}))).
		Complete(r)
}
