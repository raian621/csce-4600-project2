package builtins

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

var AliasMap map[string]string = make(map[string]string, 0)

// https://www.gnu.org/software/bash/manual/html_node/Bash-Builtins.html#index-alias
func Alias(w io.Writer, args ...string) error {
	printAllAliases := func() {
		keys := make([]string, 0, len(AliasMap))
		for k := range AliasMap {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		for _, key := range keys {
			fmt.Fprintf(w, "alias %s='%s'\n", key, AliasMap[key])
		}
	}

	if len(args) == 0 {
		printAllAliases()
	}

	for i, arg := range args {
		if i == 0 && arg == "-p" {
			printAllAliases()
		} else if j := strings.IndexByte(arg, '='); j != -1 {
			name := arg[:j]
			value := arg[j+1:]
			// remove quotes
			if value[0] == '"' || value[0] == '\'' || value[len(value)-1] == '"' || value[len(value)-1] == '\'' {
				value = value[1 : len(value)-1]
			}
			AliasMap[name] = value
		} else if value, ok := AliasMap[arg]; ok {
			fmt.Fprintf(w, "alias %s='%s'\n", arg, value)
		} else {
			fmt.Fprintf(w, "alias '%s' not found\n", arg)
		}
	}

	return nil
}
