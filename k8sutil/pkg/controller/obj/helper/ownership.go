package helper

import (
	"context"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Own(ctx context.Context, target client.Object, owner lifecycle.TypedObject) error {
	ownerAccessor, err := meta.Accessor(owner.Object)
	if err != nil {
		return err
	}

	if target.GetNamespace() != ownerAccessor.GetNamespace() {
		return &OwnerInOtherNamespaceError{
			Owner:     owner,
			OwnerKey:  client.ObjectKey{Namespace: ownerAccessor.GetNamespace(), Name: ownerAccessor.GetName()},
			Target:    target,
			TargetKey: client.ObjectKey{Namespace: target.GetNamespace(), Name: target.GetName()},
		}
	}

	if labelValue, found := ManagedByLabelValueFromContext(ctx); found && labelValue != "" {
		Label(target, "app.kubernetes.io/managed-by", labelValue)
	}

	ref := metav1.NewControllerRef(ownerAccessor, owner.GVK)

	targetOwners := target.GetOwnerReferences()
	for i, c := range targetOwners {
		if equality.Semantic.DeepEqual(c, *ref) {
			return nil
		} else if c.Controller != nil && *c.Controller {
			c.Controller = func(b bool) *bool { return &b }(false)
			klog.Warningf(
				"%T %s/%s is stealing controller for %T %s/%s from %s %s/%s",
				owner.Object, ownerAccessor.GetNamespace(), ownerAccessor.GetName(),
				target, target.GetNamespace(), target.GetName(),
				c.Kind, target.GetNamespace(), c.Name,
			)

			targetOwners[i] = c
		}
	}

	targetOwners = append(targetOwners, *ref)
	target.SetOwnerReferences(targetOwners)

	return nil
}
