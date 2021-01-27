package helper

import (
	"context"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CreateOrUpdateOptions struct {
	ObjectKey client.ObjectKey
}

type CreateOrUpdateOption interface {
	ApplyToCreateOrUpdateOptions(target *CreateOrUpdateOptions)
}

func (o *CreateOrUpdateOptions) ApplyOptions(opts []CreateOrUpdateOption) {
	for _, opt := range opts {
		opt.ApplyToCreateOrUpdateOptions(o)
	}
}

func CreateOrUpdate(ctx context.Context, cl client.Client, obj client.Object, opts ...CreateOrUpdateOption) error {
	o := &CreateOrUpdateOptions{
		ObjectKey: client.ObjectKeyFromObject(obj),
	}
	o.ApplyOptions(opts)

	obj.SetNamespace(o.ObjectKey.Namespace)
	obj.SetName(o.ObjectKey.Name)

	if Exists(obj) {
		klog.Infof("updating %T %s", obj, client.ObjectKeyFromObject(obj))
		return cl.Update(ctx, obj)
	}

	if obj.GetName() != "" {
		klog.Infof("creating %T %s", obj, client.ObjectKeyFromObject(obj))
	} else if obj.GetGenerateName() != "" {
		klog.Infof("creating %T %s%s (server-generated name)", obj, client.ObjectKeyFromObject(obj), obj.GetGenerateName())
	}
	return cl.Create(ctx, obj)
}

type PatchOptions struct {
	ObjectKey client.ObjectKey
}

type PatchOption interface {
	ApplyToPatchOptions(target *PatchOptions)
}

func (o *PatchOptions) ApplyOptions(opts []PatchOption) {
	for _, opt := range opts {
		opt.ApplyToPatchOptions(o)
	}
}

func Patch(ctx context.Context, cl client.Client, upd, orig client.Object, opts ...PatchOption) error {
	o := &PatchOptions{
		ObjectKey: client.ObjectKeyFromObject(upd),
	}
	o.ApplyOptions(opts)

	upd.SetNamespace(o.ObjectKey.Namespace)
	upd.SetName(o.ObjectKey.Name)

	klog.Infof("patching %T %s", upd, client.ObjectKeyFromObject(upd))
	return cl.Patch(ctx, upd, client.MergeFromWithOptions(orig, client.MergeFromWithOptimisticLock{}))
}

type Patcher struct {
	upd, orig client.Object
	opts      []PatchOption
}

var _ lifecycle.Persister = &Patcher{}

func (p *Patcher) Persist(ctx context.Context, cl client.Client) error {
	return Patch(ctx, cl, p.upd, p.orig, p.opts...)
}

func NewPatcher(upd, orig client.Object, opts ...PatchOption) *Patcher {
	return &Patcher{
		upd:  upd,
		orig: orig,
		opts: opts,
	}
}

func Exists(obj client.Object) bool {
	return len(obj.GetUID()) > 0
}

func GetIgnoreNotFound(ctx context.Context, cl client.Client, key client.ObjectKey, obj client.Object) (bool, error) {
	if err := cl.Get(ctx, key, obj); errors.IsNotFound(err) {
		klog.V(2).Infof("object %T %s not found", obj, key)

		obj.SetNamespace(key.Namespace)
		obj.SetName(key.Name)

		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteIgnoreNotFound(ctx context.Context, cl client.Client, obj client.Object, opts ...lifecycle.DeleteOption) (bool, error) {
	if !Exists(obj) {
		return false, nil
	}

	o := &lifecycle.DeleteOptions{}
	o.ApplyOptions(opts)

	var copts []client.DeleteOption
	if o.PropagationPolicy != "" {
		copts = append(copts, client.PropagationPolicy(o.PropagationPolicy))
	}

	klog.Infof("deleting %T %s", obj, client.ObjectKeyFromObject(obj))
	if err := cl.Delete(ctx, obj, copts...); errors.IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

type WithObjectKey client.ObjectKey

var _ CreateOrUpdateOption = WithObjectKey(client.ObjectKey{})
var _ PatchOption = WithObjectKey(client.ObjectKey{})

func (wok WithObjectKey) ApplyToCreateOrUpdateOptions(opts *CreateOrUpdateOptions) {
	opts.ObjectKey = client.ObjectKey(wok)
}

func (wok WithObjectKey) ApplyToPatchOptions(opts *PatchOptions) {
	opts.ObjectKey = client.ObjectKey(wok)
}
