package exec

import (
	"context"

	corev1obj "github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/api/corev1"
	"k8s.io/client-go/rest"
)

// ShellScriptOptions customizes shell script execution.
type ShellScriptOptions struct {
	CommandOptions

	// Shell is the path to the shell to use for executing the script. Defaults
	// to /bin/sh.
	Shell string
}

// ShellScriptOption is a setter for one or more shell script execution options.
type ShellScriptOption interface {
	ApplyToShellScriptOptions(target *ShellScriptOptions)
}

var _ ShellScriptOption = ShellScriptOptions{}

// ApplyOptions runs each of the given options against this shell script
// execution options struct.
func (o *ShellScriptOptions) ApplyOptions(opts []ShellScriptOption) {
	for _, opt := range opts {
		opt.ApplyToShellScriptOptions(o)
	}
}

// ApplyToShellScriptOptions copies all of the options set on this struct to
// another set of shell script execution options.
func (o ShellScriptOptions) ApplyToShellScriptOptions(target *ShellScriptOptions) {
	if o.Shell != "" {
		target.Shell = o.Shell
	}

	o.CommandOptions.ApplyToCommandOptions(&target.CommandOptions)
}

// ShellScript executes an arbitrary script in the specified pod.
//
// Other than presenting a simplified way of specifying the content of the
// script to run, it behaves identically to Command.
func ShellScript(ctx context.Context, cfg *rest.Config, pod *corev1obj.Pod, script string, opts ...ShellScriptOption) (*Result, error) {
	o := &ShellScriptOptions{
		Shell: "/bin/sh",
	}
	o.ApplyOptions(opts)

	return Command(ctx, cfg, pod, []string{o.Shell, "-c", script}, o.CommandOptions)
}
