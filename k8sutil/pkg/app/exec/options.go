package exec

// WithContainer specifies a container name to use for command execution.
type WithContainer string

var _ CommandOption = WithContainer("")

// ApplyToShellScriptOptions sets the container name when executing a shell
// script.
func (wc WithContainer) ApplyToShellScriptOptions(target *ShellScriptOptions) {
	wc.ApplyToCommandOptions(&target.CommandOptions)
}

// ApplyToCommandOptions sets the container name when executing a command
// directly.
func (wc WithContainer) ApplyToCommandOptions(target *CommandOptions) {
	target.Container = string(wc)
}

// WithShell specifies the name of the POSIX-compatible shell to use when
// executing a shell script.
type WithShell string

var _ ShellScriptOption = WithShell("")

// ApplyToShellScriptOptions sets the shell when executing a shell script.
func (ws WithShell) ApplyToShellScriptOptions(target *ShellScriptOptions) {
	target.Shell = string(ws)
}
