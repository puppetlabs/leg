package batchv1

import (
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/helper"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	JobKind = batchv1.SchemeGroupVersion.WithKind("Job")
)

type Job struct {
	*helper.NamespaceScopedAPIObject

	Key    client.ObjectKey
	Object *batchv1.Job
}

func makeJob(key client.ObjectKey, obj *batchv1.Job) *Job {
	j := &Job{Key: key, Object: obj}
	j.NamespaceScopedAPIObject = helper.ForNamespaceScopedAPIObject(&j.Key, lifecycle.TypedObject{GVK: JobKind, Object: j.Object})
	return j
}

func (j *Job) Copy() *Job {
	return makeJob(j.Key, j.Object.DeepCopy())
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
	return makeJob(key, &batchv1.Job{})
}

func NewJobFromObject(obj *batchv1.Job) *Job {
	return makeJob(client.ObjectKeyFromObject(obj), obj)
}

func NewJobPatcher(upd, orig *Job) lifecycle.Persister {
	return helper.NewPatcher(upd.Object, orig.Object, helper.WithObjectKey(upd.Key))
}
