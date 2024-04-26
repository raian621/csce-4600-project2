package builtins

import "io"

var AliasMap map[string]string = make(map[string]string, 0)

// https://www.gnu.org/software/bash/manual/html_node/Bash-Builtins.html#index-alias
func Alias(w io.Writer, args ...string) error {
	return nil
}
