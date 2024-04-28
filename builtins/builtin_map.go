package builtins

import "io"

type BuiltinEntrypoint func(io.Writer, ...string) error

var BuiltinMap map[string]BuiltinEntrypoint

func init() {
	BuiltinMap = map[string]BuiltinEntrypoint{
		"cd": func(_w io.Writer, args ...string) error {
			return ChangeDirectory(args...)
		},
		"env": EnvironmentVariables,

		// Ryan-implemented builtins
		"alias":   Alias,
		"builtin": Builtin,
		"history": History,
		"pwd":     PrintWorkingDirectory,
		"unalias": Unalias,
	}
}
