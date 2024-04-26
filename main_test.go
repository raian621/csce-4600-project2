package main

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"
	"time"

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

func Test_tokenizeInput(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    []string
		wantErr error
	}{
		{
			name:  "basic tokenization",
			input: "ls Homework",
			want: []string{
				"ls", "Homework",
			},
		},
		{
			name:  "adversarial whitespace",
			input: "    ls   Homework  ",
			want: []string{
				"ls", "Homework",
			},
		},
		{
			name:  "input with \" string with spaces",
			input: "ls \"Definitely a Homework Folder\"",
			want: []string{
				"ls", "Definitely a Homework Folder",
			},
		},
		{
			name:  "input with ' string with spaces",
			input: "ls 'Definitely a Homework Folder'",
			want: []string{
				"ls", "Definitely a Homework Folder",
			},
		},
		{
			name:  "input with escaped characters",
			input: "ls \\'Definitely a Homework Folder\\'",
			want: []string{
				"ls", "\\'Definitely", "a", "Homework", "Folder\\'",
			},
		},
		{
			name:    "no closing quote \"",
			input:   "ls \"Definitely a Homework Folder",
			wantErr: ErrNoClosingQuote,
		},
		{
			name:  "valid nested quotes",
			input: "ls \"Definitely a 'Homework' Folder\"",
			want: []string{
				"ls", "Definitely a 'Homework' Folder",
			},
		},
		{
			name:  "setting an arg equal to a string",
			input: "alias uatemycookie=\"low tc\"",
			want: []string{
				"alias", "uatemycookie=\"low tc\"",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // avoid Go loop gotcha ig

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tokens, err := tokenizeInput(tc.input)
			// reflect.DeepEqual was returning false when comparing empty arrays
			// for some reason:
			if !reflect.DeepEqual(tokens, tc.want) && (len(tokens) != 0 || len(tc.want) != 0) {
				t.Errorf("wanted %v, got %v tokens", tc.want, tokens)
			}
			if err != tc.wantErr {
				t.Errorf("wanted %v, got %v error", tc.wantErr, err)
			}
		})
	}
}
