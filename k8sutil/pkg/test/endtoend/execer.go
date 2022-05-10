package endtoend

import (
	"context"
	"fmt"
	"time"

	"github.com/puppetlabs/leg/k8sutil/pkg/app/exec"
	corev1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/corev1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ExecerOptions allow for customization of a script executor.
type ExecerOptions struct {
	// PodMeta is the pod metadata for the execution environment. Defaults to
	// the default namespace and a generated name starting with "script-".
	PodMeta metav1.ObjectMeta

	// Image is the name of the Docker image to use. Defaults to Alpine.
	Image string

	// Shell is the POSIX-compatible shell to use. Defaults to /bin/sh.
	Shell string

	// Timeout is the maximum lifetime for the script executor pod. Defaults to
	// 24 hours.
	Timeout time.Duration
}

// ExecerOption is a setter for one or more script executor options.
type ExecerOption interface {
	// ApplyToExecerOptions copies the configuration of this option to the given
	// script executor options.
	ApplyToExecerOptions(target *ExecerOptions)
}

// ApplyOptions runs each of the given options against this script executor
// options struct.
func (o *ExecerOptions) ApplyOptions(opts []ExecerOption) {
	for _, opt := range opts {
		opt.ApplyToExecerOptions(o)
	}
}

// ExecerWithNamespace causes the script execution pod to be created in the
// given namespace.
type ExecerWithNamespace string

var _ ExecerOption = ExecerWithNamespace("")

// ApplyToExecerOptions copies the configuration of this option to the given
// script executor options.
func (ewn ExecerWithNamespace) ApplyToExecerOptions(target *ExecerOptions) {
	target.PodMeta.SetNamespace(string(ewn))
}

// ExecerWithName sets the pod name.
type ExecerWithName string

var _ ExecerOption = ExecerWithName("")

// ApplyToExecerOptions copies the configuration of this option to the given
// script executor options.
func (ewn ExecerWithName) ApplyToExecerOptions(target *ExecerOptions) {
	target.PodMeta.SetGenerateName("")
	target.PodMeta.SetName(string(ewn))
}

// ExecerWithGenerateName generates the pod name from the given template string.
type ExecerWithGenerateName string

var _ ExecerOption = ExecerWithGenerateName("")

// ApplyToExecerOptions copies the configuration of this option to the given
// script executor options.
func (ewgn ExecerWithGenerateName) ApplyToExecerOptions(target *ExecerOptions) {
	target.PodMeta.SetGenerateName(string(ewgn))
	target.PodMeta.SetName("")
}

// ExecerWithImage sets the Docker image to use for the execution environment.
type ExecerWithImage string

var _ ExecerOption = ExecerWithImage("")

// ApplyToExecerOptions copies the configuration of this option to the given
// script executor options.
func (ewi ExecerWithImage) ApplyToExecerOptions(target *ExecerOptions) {
	target.Image = string(ewi)
}

// ExecerWithShell sets the POSIX-compatible shell to use when executing
// scripts.
type ExecerWithShell string

var _ ExecerOption = ExecerWithShell("")

// ApplyToExecerOptions copies the configuration of this option to the given
// script executor options.
func (ews ExecerWithShell) ApplyToExecerOptions(target *ExecerOptions) {
	target.Shell = string(ews)
}

// ExecerWithTimeout sets the maximum lifetime for the pod.
type ExecerWithTimeout time.Duration

var _ ExecerOption = ExecerWithTimeout(0)

// ApplyToExecerOptions copies the configuration of this option to the given
// script executor options.
func (ewt ExecerWithTimeout) ApplyToExecerOptions(target *ExecerOptions) {
	target.Timeout = time.Duration(ewt)
}

// Execer is a utility that allows arbitrary shell commands to be run inside a
// cluster.
//
// For example, this can be used to test that services only accessible within a
// cluster are behaving correctly without port-forwarding them.
type Execer struct {
	e     *Environment
	pod   *corev1obj.Pod
	shell string
}

// Close tears down the pod being used to execute scripts.
func (e *Execer) Close(ctx context.Context) (err error) {
	_, err = e.pod.Delete(ctx, e.e.ControllerClient)
	return
}

// Exec executes the given script using a POSIX-compatible shell.
//
// The result of executing the shell is made available in the result struct. If
// the command fails, its exit status is reported in the result, but an error is
// not returned. This method only returns an error if an infrastucture failure
// occurs (like not being able to communicate with the execution pod).
func (e *Execer) Exec(ctx context.Context, script string) (*exec.Result, error) {
	if err := e.pod.Persist(ctx, e.e.ControllerClient); err != nil {
		return nil, err
	}

	return exec.ShellScript(ctx, e.e.RESTConfig, e.pod, script, exec.WithShell(e.shell))
}

// NewExecer creates a new script executor for the given environment and
// options.
func NewExecer(e *Environment, opts ...ExecerOption) *Execer {
	o := &ExecerOptions{
		PodMeta: metav1.ObjectMeta{
			Namespace:    "default",
			GenerateName: "script-",
		},
		Image:   "alpine:latest",
		Shell:   "/bin/sh",
		Timeout: 24 * time.Hour,
	}
	o.ApplyOptions(opts)

	pod := &corev1.Pod{
		ObjectMeta: o.PodMeta,
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "script",
					Image: o.Image,
					Args: []string{
						o.Shell,
						"-c",
						fmt.Sprintf("trap : TERM INT; sleep %d & wait", o.Timeout/time.Second),
					},
				},
			},
		},
	}
	return &Execer{
		e:     e,
		pod:   corev1obj.NewPodFromObject(pod),
		shell: o.Shell,
	}
}

// Exec creates a one-time-use Execer, runs the given script, and then tears
// down the backing pod.
func Exec(ctx context.Context, e *Environment, script string, opts ...ExecerOption) (res *exec.Result, err error) {
	execer := NewExecer(e, opts...)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if cerr := execer.Close(ctx); cerr != nil {
			if err == nil {
				err = cerr
			}
		}
	}()

	return execer.Exec(ctx, script)
}
