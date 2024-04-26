package builtins

import "io"

// https://www.gnu.org/software/bash/manual/html_node/Bash-Builtins.html#index-unalias
func Unalias(w io.Writer, args ...string) error {
	return nil
}
