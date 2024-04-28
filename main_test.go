package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"testing/iotest"
	"time"

	"github.com/raian621/csce-4600-project2/builtins"
	"github.com/stretchr/testify/require"
)

func Test_runLoop(t *testing.T) {
	t.Parallel()
	exitCmd := strings.NewReader("exit\n")
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name     string
		args     args
		wantW    string
		wantErrW string
	}{
		{
			name: "no error",
			args: args{
				r: exitCmd,
			},
		},
		{
			name: "read error should have no effect",
			args: args{
				r: iotest.ErrReader(io.EOF),
			},
			wantErrW: "EOF",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := &bytes.Buffer{}
			errW := &bytes.Buffer{}

			exit := make(chan struct{}, 2)
			// run the loop for 10ms
			go runLoop(tt.args.r, w, errW, exit)
			time.Sleep(10 * time.Millisecond)
			exit <- struct{}{}

			require.NotEmpty(t, w.String())
			if tt.wantErrW != "" {
				require.Contains(t, errW.String(), tt.wantErrW)
			} else {
				require.Empty(t, errW.String())
			}
		})
	}
}

func Test_executeCommand(t *testing.T) {
	args := []string{"echo", "ligma"}
	err := executeCommand(args[0], args[1:]...)

	require.Nil(t, err)
}

func TestHandleEmptyInput(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	err := handleInput(&out, "", make(chan<- struct{}))

	require.Nil(t, err)
	require.Equal(t, "", out.String())
}

func TestAliasSystem(t *testing.T) {
	builtins.AliasMap = map[string]string{
		"a": "alias",
	}

	var out bytes.Buffer
	err := handleInput(&out, "a", make(chan<- struct{}))

	require.Nil(t, err)
	require.Equal(t, "alias a='alias'\n", out.String())
}

func TestTokenizerErrHandling(t *testing.T) {
	var out bytes.Buffer
	err := handleInput(&out, "ls -l \"", make(chan<- struct{}))

	require.Error(t, ErrNoClosingQuote, err)
	require.Equal(t, "", out.String())
}

func TestTokenizerWithAliasErrHandling(t *testing.T) {
	var out bytes.Buffer
	err := handleInput(&out, "ls -l \"", make(chan<- struct{}))

	require.Error(t, ErrNoClosingQuote, err)
	require.Equal(t, "", out.String())
}
