package exec

import (
	"bytes"
	"context"
	"net/http"

	corev1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/corev1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/exec"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CommandOptions customizes command execution.
type CommandOptions struct {
	// Container is the name of the container to run the command in. If not
	// specified, defaults to the first container in the pod.
	Container string
}

// CommandOption is a setter for one or more command or script execution
// options.
type CommandOption interface {
	ShellScriptOption

	// ApplyToCommandOptions copies the configuration of this option to the
	// given command execution options.
	ApplyToCommandOptions(target *CommandOptions)
}

var _ CommandOption = CommandOptions{}

// ApplyOptions runs each of the given options against this command execution
// options struct.
func (o *CommandOptions) ApplyOptions(opts []CommandOption) {
	for _, opt := range opts {
		opt.ApplyToCommandOptions(o)
	}
}

// ApplyToShellScriptOptions copies all of the options set on this struct to the
// given shell script execution options.
func (o CommandOptions) ApplyToShellScriptOptions(target *ShellScriptOptions) {
	o.ApplyToCommandOptions(&target.CommandOptions)
}

// ApplyToCommandOptions copies all of the options set on this struct to another
// set of command execution options.
func (o CommandOptions) ApplyToCommandOptions(target *CommandOptions) {
	if o.Container != "" {
		target.Container = o.Container
	}
}

// Result provides access to some information after running a command.
type Result struct {
	// ExitCode is the exit code returned by the command.
	ExitCode int

	// Stdout is the complete contents of the command's standard output file
	// descriptor.
	Stdout string

	// Stderr is the complete contents of the command's standard error file
	// descriptor.
	Stderr string
}

// Command executes an arbitrary command in the specified pod.
//
// The result of executing the command is provided in the returned struct. If
// the command fails, this function will not return an error, but the result
// ExitCode will be set to a non-zero value.
func Command(ctx context.Context, cfg *rest.Config, pod *corev1obj.Pod, command []string, opts ...CommandOption) (*Result, error) {
	o := &CommandOptions{}
	o.ApplyOptions(opts)

	cl, err := client.New(cfg, client.Options{})
	if err != nil {
		return nil, err
	}

	kc, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	if _, err := corev1obj.NewPodRunningPoller(pod).Load(ctx, cl); err != nil {
		return nil, err
	}

	if o.Container == "" {
		o.Container = pod.Object.Spec.Containers[0].Name
	}

	if _, err := corev1obj.NewPodContainerRunningPoller(pod, o.Container).Load(ctx, cl); err != nil {
		return nil, err
	}

	req := kc.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(pod.Key.Namespace).
		Name(pod.Key.Name).
		SubResource("exec").
		Param("container", o.Container)
	req.VersionedParams(&corev1.PodExecOptions{
		Container: o.Container,
		Command:   command,
		Stdout:    true,
		Stderr:    true,
	}, scheme.ParameterCodec)

	execer, err := remotecommand.NewSPDYExecutor(cfg, http.MethodPost, req.URL())
	if err != nil {
		return nil, err
	}

	var stdout, stderr bytes.Buffer
	err = execer.Stream(remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	})

	var code int
	if cerr, ok := err.(exec.CodeExitError); ok {
		code = cerr.ExitStatus()
	} else if err != nil {
		return nil, err
	}

	return &Result{
		ExitCode: code,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
	}, nil
}
