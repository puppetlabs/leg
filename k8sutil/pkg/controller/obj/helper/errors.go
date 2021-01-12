package helper

import (
	"fmt"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type OwnerInOtherNamespaceError struct {
	Owner  lifecycle.TypedObject
	Target client.Object
}

func (e *OwnerInOtherNamespaceError) Error() string {
	return fmt.Sprintf(
		"owner %T %s is in a different namespace than %T %s",
		e.Owner.Object,
		client.ObjectKeyFromObject(e.Owner.Object),
		e.Target,
		client.ObjectKeyFromObject(e.Target),
	)
}
