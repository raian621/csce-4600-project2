package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/raian621/csce-4600-project2/builtins"
)

type BuiltinEntrypoint func(io.Writer, ...string) error

var builtinMap map[string]BuiltinEntrypoint = map[string]BuiltinEntrypoint{
	"cd": func(_w io.Writer, args ...string) error {
		return builtins.ChangeDirectory(args...)
	},
	"env": builtins.EnvironmentVariables,

	// Ryan-implemented builtins
	"alias":   builtins.Alias,
	"unalias": builtins.Unalias,
	"pwd":     builtins.PrintWorkingDirectory,
}

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
	if fn, ok := builtinMap[name]; ok {
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

var ErrNoClosingQuote = errors.New("missing closing quote")

func tokenizeInput(input string) ([]string, error) {
	var (
		start        int      = 0
		end          int      = 0
		tokens       []string = make([]string, 0)
		shouldEscape bool     = false
	)

	// om nom nom
	consumeToken := func(start, end int) {
		tokens = append(tokens, input[start:end])
	}

	for end < len(input) {
		if shouldEscape {
			end++
			shouldEscape = false
			continue
		}

		switch input[end] {
		// spaces separate tokens
		case ' ':
			if start != end {
				consumeToken(start, end)
			}
			start = end + 1
		// escape certain characters
		case '\\':
			shouldEscape = true
		// consume string token
		case '\'':
			fallthrough
		case '"':
			quoteChar := input[end]
			end++
			for end < len(input) && input[end] != quoteChar {
				end++
			}
			if end == len(input) {
				return []string{}, ErrNoClosingQuote
			}
			if input[start] != quoteChar {
				// this case is for something like cmd var="askdjfhkdfhj"
				consumeToken(start, end+1)
			} else {
				// do not consume leading quote
				consumeToken(start+1, end)
			}
			start = end + 1
		}
		end++
	}

	// consume last token if it exists
	if start != end {
		consumeToken(start, end)
	}

	return tokens, nil
}
