package batchv1

import (
	"context"

	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
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

func (j *Job) Delete(ctx context.Context, cl client.Client, opts ...lifecycle.DeleteOption) (bool, error) {
	return helper.DeleteIgnoreNotFound(ctx, cl, j.Object, opts...)
}

func (j *Job) LabelAnnotateFrom(ctx context.Context, from metav1.Object) {
	helper.CopyLabelsAndAnnotations(&j.Object.ObjectMeta, from)
}

func (j *Job) Load(ctx context.Context, cl client.Client) (bool, error) {
	return helper.GetIgnoreNotFound(ctx, cl, j.Key, j.Object)
}

func (j *Job) Owned(ctx context.Context, owner lifecycle.TypedObject) error {
	return helper.Own(j.Object, owner)
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

func (j *Job) Condition(typ batchv1.JobConditionType) (batchv1.JobCondition, bool) {
	for _, cond := range j.Object.Status.Conditions {
		if cond.Type == typ {
			return cond, true
		}
	}
	return batchv1.JobCondition{Type: typ}, false
}

func (j *Job) CompleteCondition() (batchv1.JobCondition, bool) {
	return j.Condition(batchv1.JobComplete)
}

func (j *Job) FailedCondition() (batchv1.JobCondition, bool) {
	return j.Condition(batchv1.JobFailed)
}

func (j *Job) Complete() bool {
	cc, ok := j.CompleteCondition()
	return ok && cc.Status == corev1.ConditionTrue
}

func (j *Job) Failed() bool {
	fc, ok := j.FailedCondition()
	return ok && fc.Status == corev1.ConditionTrue
}

func (j *Job) Succeeded() bool {
	return j.Complete() && !j.Failed()
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
