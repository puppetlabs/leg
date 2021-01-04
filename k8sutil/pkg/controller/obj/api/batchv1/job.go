package batchv1

import (
	"context"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	JobKind = batchv1.SchemeGroupVersion.WithKind("Job")
)

type Job struct {
	Key    client.ObjectKey
	Object *batchv1.Job
}

var _ lifecycle.Deleter = &Job{}
var _ lifecycle.LabelAnnotatableFrom = &Job{}
var _ lifecycle.Loader = &Job{}
var _ lifecycle.Ownable = &Job{}
var _ lifecycle.Persister = &Job{}

func (j *Job) Delete(ctx context.Context, cl client.Client) (bool, error) {
	return helper.DeleteIgnoreNotFound(ctx, cl, j.Object)
}

func (j *Job) LabelAnnotateFrom(ctx context.Context, from metav1.Object) {
	helper.CopyLabelsAndAnnotations(&j.Object.ObjectMeta, from)
}

func (j *Job) Load(ctx context.Context, cl client.Client) (bool, error) {
	return helper.GetIgnoreNotFound(ctx, cl, j.Key, j.Object)
}

func (j *Job) Owned(ctx context.Context, owner lifecycle.TypedObject) error {
	return helper.Own(ctx, j.Object, owner)
}

func (j *Job) Persist(ctx context.Context, cl client.Client) error {
	return helper.CreateOrUpdate(ctx, cl, j.Object, helper.WithObjectKey(j.Key))
}

func (j *Job) Copy() *Job {
	return &Job{
		Key:    j.Key,
		Object: j.Object.DeepCopy(),
	}
}

func NewJob(key client.ObjectKey) *Job {
	return &Job{
		Key:    key,
		Object: &batchv1.Job{},
	}
}

func NewJobFromObject(obj *batchv1.Job) *Job {
	return &Job{
		Key:    client.ObjectKeyFromObject(obj),
		Object: obj,
	}
}

func NewJobPatcher(upd, orig *Job) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(upd.Key))
}
