package fakeext_test

import (
	"context"
	"testing"

	"github.com/puppetlabs/leg/k8sutil/pkg/test/fakeext"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestSetUIDExtension(t *testing.T) {
	ctx := context.Background()

	obj := &corev1.Namespace{}
	obj.SetName("my-first-namespace")

	kc, err := fakeext.NewKubernetesClientsetWithExtensions(
		[]runtime.Object{obj},
		[]fakeext.Extension{fakeext.SetUIDExtension},
	)
	require.NoError(t, err)

	obj, err = kc.CoreV1().Namespaces().Get(ctx, "my-first-namespace", metav1.GetOptions{})
	require.NoError(t, err)
	require.NotEmpty(t, obj.GetUID())

	obj = &corev1.Namespace{}
	obj.SetName("my-second-namespace")

	obj, err = kc.CoreV1().Namespaces().Create(ctx, obj, metav1.CreateOptions{})
	require.NoError(t, err)
	require.NotEmpty(t, obj.GetUID())
}
