package helper

import (
	"fmt"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type OwnerInOtherNamespaceError struct {
	Owner     lifecycle.TypedObject
	OwnerKey  client.ObjectKey
	Target    runtime.Object
	TargetKey client.ObjectKey
}

func (e *OwnerInOtherNamespaceError) Error() string {
	return fmt.Sprintf("owner %T %s is in a different namespace than %T %s", e.Owner.Object, e.OwnerKey, e.Target, e.TargetKey)
}
