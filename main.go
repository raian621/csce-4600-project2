package main

import (
	"bufio"
	"bytes"
	"errors"
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

	// Split the input separate the command name and the command arguments.
	args := strings.Split(input, " ")
	name, args := args[0], args[1:]

	// Check for built-in commands.
	// New builtin commands should be added here. Eventually this should be refactored to its own func.
	switch name {
	case "cd":
		return builtins.ChangeDirectory(args...)
	case "env":
		return builtins.EnvironmentVariables(w, args...)
	case "exit":
		exit <- struct{}{}
		return nil
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

func tokenizeInput(input string, removeEscapeChars bool) ([]string, error) {
	var (
		start        int      = 0
		end          int      = 0
		tokens       []string = make([]string, 0)
		shouldEscape bool     = false
	)

	filterEscapeChars := func(token string) string {
		var filtered bytes.Buffer

		var i int
		for i < len(token) {
			if token[i] == '\\' {
				i++
			}
			filtered.WriteByte(token[i])
			i++
		}

		return filtered.String()
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
				token := string(input[start:end])
				if removeEscapeChars {
					token = filterEscapeChars(token)
				}
				tokens = append(tokens, token)
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
			token := string(input[start+1 : end])
			if removeEscapeChars {
				token = filterEscapeChars(token)
			}
			tokens = append(tokens, token)
			start = end + 1
		}
		end++
	}

	// consume last token if it exists
	if start != end {
		token := string(input[start:end])
		if removeEscapeChars {
			token = filterEscapeChars(token)
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}
