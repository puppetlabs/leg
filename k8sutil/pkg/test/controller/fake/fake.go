package fake

import (
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kubernetesscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func NewClient() client.Client {
	scheme := runtime.NewScheme()
	schemeBuilder := runtime.NewSchemeBuilder(
		kubernetesscheme.AddToScheme,
		metav1.AddMetaToScheme,
		rbacv1.AddToScheme,
		apiextensionsv1.AddToScheme,
		apiextensionsv1beta1.AddToScheme,
	)
	_ = schemeBuilder.AddToScheme(scheme)

	return fake.NewClientBuilder().
		WithObjects().
		WithScheme(scheme).
		Build()
}
