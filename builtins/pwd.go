package builtins

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var ErrUnknownArg error = errors.New("unknown argument supplied")

// https://www.gnu.org/software/bash/manual/html_node/Bourne-Shell-Builtins.html#index-pwd
func PrintWorkingDirectory(w io.Writer, args ...string) (err error) {
	if len(args) > 1 {
		fmt.Fprintf(w, "usage: pwd [-L|-P]")
		return ErrInvalidArgCount
	}

	var cwd string

	if len(args) == 0 {
		// get the "physical" working directory by default
		// probably a safer bet?
		cwd, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(w, "unexpected error occurred '%v'", err)
			return err
		}
	} else {
		flag := args[0]

		if flag == "-L" {
			// get "logical" working directory
			cwd = os.Getenv("PWD")
		} else if flag == "-P" {
			// get "physical" working directory
			cwd, err = os.Getwd()
			if err != nil {
				fmt.Fprintf(w, "unexpected error occurred '%v'", err)
				return err
			}
		} else {
			fmt.Fprint(w, "usage: pwd [-L|-P]")
			return ErrUnknownArg
		}
	}

	fmt.Fprintln(w, cwd)

	return nil
}
