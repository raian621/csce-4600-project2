package builtins_test

import (
	"bytes"
	"testing"

	"github.com/raian621/csce-4600-project2/builtins"
)

func TestBuiltin(t *testing.T) {
	testCases := []struct {
		name    string
		args    []string
		wantErr error
	}{
		{
			name: "run a builtin command",
			args: []string{"builtin", "builtin"},
		},
		{
			name:    "run a builtin command",
			args:    []string{"builtin", "ls"},
			wantErr: builtins.ErrNotBuiltinCommand,
		},
	}

	for _, tc := range testCases {
		tc := tc
		var out bytes.Buffer

		t.Run(tc.name, func(t *testing.T) {
			if err := builtins.Builtin(&out, tc.args...); err != tc.wantErr {
				t.Fatalf("expected '%v', got '%v' error", tc.wantErr, err)
			}
		})
	}
}
