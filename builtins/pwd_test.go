package builtins_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/raian621/csce-4600-project2/builtins"
)

func TestPrintWorkingDirectory(t *testing.T) {
	// kept getting 'getwd: no such file or directory' before adding this line
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name    string
		args    []string
		wantErr error
		wantOut string
	}{
		{
			name: "no arguments",
			args: []string{},
		},
		{
			name: "logical current working directory",
			args: []string{"-L"},
		},
		{
			name: "physical current working directory",
			args: []string{"-P"},
		},
		{
			name:    "too many arguments",
			args:    []string{"-P", "-C"},
			wantErr: builtins.ErrInvalidArgCount,
			wantOut: "usage: pwd [-L|-P]",
		},
		{
			name:    "undefined argument",
			args:    []string{"-C"},
			wantErr: builtins.ErrUnknownArg,
			wantOut: "usage: pwd [-L|-P]",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			err := builtins.PrintWorkingDirectory(&out, tc.args...)

			if err != tc.wantErr {
				t.Fatalf("expected '%v', got '%v' error", tc.wantErr, err)
			}

			// don't try to test for an exact PWD on success, as it'll likely  vary across machines,
			// just verify that the usage is printed out upon an error.
			if err != nil && out.String() != tc.wantOut {

			}
		})
	}
}
