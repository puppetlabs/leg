package manifest_test

import (
	"testing"

	"github.com/puppetlabs/leg/k8sutil/pkg/manifest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta/testrestmapper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

func TestFixupPatcher(t *testing.T) {
	deployment := &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "alpine:latest",
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 80,
								},
								{
									Name:          "dns",
									Protocol:      corev1.ProtocolUDP,
									ContainerPort: 53,
								},
							},
						},
					},
				},
			},
		},
	}
	deploymentGVK, err := apiutil.GVKForObject(deployment, scheme.Scheme)
	require.NoError(t, err)

	manifest.FixupPatcher(deployment, &deploymentGVK)
	assert.Equal(t, corev1.ProtocolTCP, deployment.Spec.Template.Spec.Containers[0].Ports[0].Protocol)
	assert.Equal(t, corev1.ProtocolUDP, deployment.Spec.Template.Spec.Containers[0].Ports[1].Protocol)

	svc := &corev1.Service{
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromString("http"),
				},
				{
					Name:       "dns",
					Protocol:   corev1.ProtocolUDP,
					Port:       53,
					TargetPort: intstr.FromString("dns"),
				},
			},
		},
	}
	svcGVK, err := apiutil.GVKForObject(svc, scheme.Scheme)
	require.NoError(t, err)

	manifest.FixupPatcher(svc, &svcGVK)
	assert.Equal(t, corev1.ProtocolTCP, svc.Spec.Ports[0].Protocol)
	assert.Equal(t, corev1.ProtocolUDP, svc.Spec.Ports[1].Protocol)
}

func TestDefaultNamespacePatcher(t *testing.T) {
	mapper := testrestmapper.TestOnlyStaticRESTMapper(scheme.Scheme)
	patcher := manifest.DefaultNamespacePatcher(mapper, "foo")

	tests := []struct {
		Name              string
		Object            client.Object
		ExpectedNamespace string
	}{
		{
			Name: "Cluster-scoped",
			Object: &corev1.PersistentVolume{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
			},
			ExpectedNamespace: "",
		},
		{
			Name: "Namespace specified",
			Object: &corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "bar",
					Name:      "test",
				},
			},
			ExpectedNamespace: "bar",
		},
		{
			Name: "Namespace not specified",
			Object: &corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
			},
			ExpectedNamespace: "foo",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			gvk, err := apiutil.GVKForObject(test.Object, scheme.Scheme)
			require.NoError(t, err)

			patcher(test.Object, &gvk)
			assert.Equal(t, test.ExpectedNamespace, test.Object.GetNamespace())
		})
	}
}
