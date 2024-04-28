package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/raian621/csce-4600-project2/builtins"
)

func main() {
	exit := make(chan struct{}, 2) // buffer this so there's no deadlock.
	runLoop(os.Stdin, os.Stdout, os.Stderr, exit)
}

func runLoop(r io.Reader, w, errW io.Writer, exit chan struct{}) {
	var (
		input    string
		err      error
		readLoop = bufio.NewReader(r)
	)
	for {
		select {
		case <-exit:
			_, _ = fmt.Fprintln(w, "exiting gracefully...")
			return
		default:
			if err := printPrompt(w); err != nil {
				_, _ = fmt.Fprintln(errW, err)
				continue
			}
			if input, err = readLoop.ReadString('\n'); err != nil {
				_, _ = fmt.Fprintln(errW, err)
				continue
			}
			if err = handleInput(w, input, exit); err != nil {
				_, _ = fmt.Fprintln(errW, err)
			}
		}
	}
}

func printPrompt(w io.Writer) error {
	// Get current user.
	// Don't prematurely memoize this because it might change due to `su`?
	u, err := user.Current()
	if err != nil {
		return err
	}
	// Get current working directory.
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// /home/User [Username] $
	_, err = fmt.Fprintf(w, "%v [%v] $ ", wd, u.Username)

	return err
}

func handleInput(w io.Writer, input string, exit chan<- struct{}) error {
	// Remove trailing spaces.
	input = strings.TrimSpace(input)

	if len(input) > 0 {
		builtins.AddHistoryEntry(input)
	} else {
		return nil
	}

	// seperate arguments in the input by tokenizing the input
	args, err := tokenizeInput(input)
	if err != nil {
		return err
	}
	name, args := args[0], args[1:]
	// if the name of the command matches a defined alias, replace the name and args[0:len(aliasTokens)]
	// with tokens generated from the defined alias:
	if alias, ok := builtins.AliasMap[name]; ok {
		aliasArgs, err := tokenizeInput(alias)
		if err != nil {
			return err
		}

		name = aliasArgs[0]
		args = append(aliasArgs[1:], args...)
	}

	if name == "exit" {
		exit <- struct{}{}
		return nil
	}

	// Check for built-in commands.
	if fn, ok := builtins.BuiltinMap[name]; ok {
		return fn(w, args...)
	}

	return executeCommand(name, args...)
}

func executeCommand(name string, arg ...string) error {
	// Otherwise prep the command
	cmd := exec.Command(name, arg...)

	// Set the correct output device.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Execute the command and return the error.
	return cmd.Run()
}
