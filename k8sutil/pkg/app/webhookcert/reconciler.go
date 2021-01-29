package webhookcert

import (
	"context"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	admissionregistrationv1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/admissionregistrationv1"
	corev1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/corev1"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups=admissionregistration.k8s.io,resources=validatingwebhookconfigurations;mutatingwebhookconfigurations,verbs=update

type Reconciler struct {
	cl                                  client.Client
	validatingWebhookConfigurationNames []string
	mutatingWebhookConfigurationNames   []string
}

func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (res reconcile.Result, err error) {
	klog.InfoS("webhook certificate reconciler: starting reconcile", "secret", req.NamespacedName)
	defer klog.InfoS("webhook certificate reconciler: ending reconcile", "secret", req.NamespacedName)
	defer func() {
		if err != nil {
			klog.ErrorS(err, "webhook certificate reconciler: failed to reconcile", "secret", req.NamespacedName)
		}
	}()

	secret := corev1obj.NewTLSSecret(req.NamespacedName)
	if _, err := (lifecycle.RequiredLoader{Loader: secret}).Load(ctx, r.cl); err != nil {
		return reconcile.Result{}, err
	}

	cert, err := secret.Certificate()
	if err != nil {
		return reconcile.Result{}, err
	}

	if len(cert.Certificate) < 2 {
		return reconcile.Result{}, fmt.Errorf("certificate in secret is missing chain")
	}

	var caBundle []byte
	for _, certDER := range cert.Certificate[1:] {
		certPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certDER,
		})
		caBundle = append(caBundle, certPEM...)
	}

	encodedCABundle := []byte(base64.StdEncoding.EncodeToString(caBundle))

	var requeue bool

	for _, name := range r.validatingWebhookConfigurationNames {
		vwc := admissionregistrationv1obj.NewValidatingWebhookConfiguration(name)
		if _, err := (lifecycle.RequiredLoader{Loader: vwc}).Load(ctx, r.cl); err != nil {
			klog.ErrorS(err, "webhook certificate reconciler: failed to load ValidatingWebhookConfiguration %q", name)
			requeue = true
			continue
		}

		for i := range vwc.Object.Webhooks {
			vwc.Object.Webhooks[i].ClientConfig.CABundle = encodedCABundle
		}

		if err := vwc.Persist(ctx, r.cl); err != nil {
			klog.ErrorS(err, "webhook certificate reconciler: failed to persist ValidatingWebhookConfiguration %q", name)
			requeue = true
			continue
		}

		klog.V(4).InfoS("webhook certificate reconciler: updated CA bundle for ValidatingWebhookConfiguration %q", name)
	}

	for _, name := range r.mutatingWebhookConfigurationNames {
		mwc := admissionregistrationv1obj.NewMutatingWebhookConfiguration(name)
		if _, err := (lifecycle.RequiredLoader{Loader: mwc}).Load(ctx, r.cl); err != nil {
			klog.ErrorS(err, "webhook certificate reconciler: failed to load MutatingWebhookConfiguration %q", name)
			requeue = true
			continue
		}

		for i := range mwc.Object.Webhooks {
			mwc.Object.Webhooks[i].ClientConfig.CABundle = encodedCABundle
		}

		if err := mwc.Persist(ctx, r.cl); err != nil {
			klog.ErrorS(err, "webhook certificate reconciler: failed to persist MutatingWebhookConfiguration %q", name)
			requeue = true
			continue
		}

		klog.V(4).InfoS("webhook certificate reconciler: updated CA bundle for MutatingWebhookConfiguration %q", name)
	}

	return reconcile.Result{Requeue: requeue}, nil
}

type Option func(r *Reconciler)

func WithValidatingWebhookConfiguration(name string) Option {
	return func(r *Reconciler) {
		r.validatingWebhookConfigurationNames = append(r.validatingWebhookConfigurationNames, name)
	}
}

func WithMutatingWebhookConfiguration(name string) Option {
	return func(r *Reconciler) {
		r.mutatingWebhookConfigurationNames = append(r.mutatingWebhookConfigurationNames, name)
	}
}

func NewReconciler(cl client.Client, opts ...Option) *Reconciler {
	r := &Reconciler{
		cl: cl,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func AddReconcilerToManager(mgr manager.Manager, secretKey client.ObjectKey, opts ...Option) error {
	r := NewReconciler(mgr.GetClient(), opts...)

	return builder.ControllerManagedBy(mgr).
		For(&corev1.Secret{}, builder.WithPredicates(predicate.NewPredicateFuncs(func(obj client.Object) bool {
			return client.ObjectKeyFromObject(obj) == secretKey
		}))).
		Complete(r)
}
