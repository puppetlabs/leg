package helper

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func SuffixObjectKeyName(key client.ObjectKey, suffix string) client.ObjectKey {
	return client.ObjectKey{
		Namespace: key.Namespace,
		Name:      fmt.Sprintf("%s-%s", key.Name, suffix),
	}
}
