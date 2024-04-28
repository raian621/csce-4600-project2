package builtins

import (
	"errors"
	"io"
)

var ErrNotBuiltinCommand error = errors.New("provided command was not a builtin command")

// https://www.gnu.org/software/bash/manual/html_node/Bash-Builtins.html#index-builtin
func Builtin(w io.Writer, args ...string) error {
	if len(args) == 0 {
		return nil
	}

	if fn, ok := BuiltinMap[args[0]]; ok {
		return fn(w, args[1:]...)
	} else {
		return ErrNotBuiltinCommand
	}
}
