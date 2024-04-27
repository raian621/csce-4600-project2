package builtins

import (
	"fmt"
	"io"
)

// https://www.gnu.org/software/bash/manual/html_node/Bash-Builtins.html#index-unalias
func Unalias(w io.Writer, args ...string) (err error) {
	if len(args) == 0 {
		fmt.Fprintln(w, "usage: unalias [-a] name [name ...]")
		return ErrInvalidArgCount
	}

	for i, arg := range args {
		if i == 0 && arg == "-a" {
			AliasMap = make(map[string]string)
		} else if _, ok := AliasMap[arg]; ok {
			delete(AliasMap, arg)
		} else {
			fmt.Fprintf(w, "alias '%s' not found\n", arg)
			err = ErrUndefinedAlias
		}
	}

	return err
}
