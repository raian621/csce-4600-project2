package builtins_test

import (
	"bytes"
	"testing"

	"github.com/raian621/csce-4600-project2/builtins"
)

func test_Alias(t *testing.T) {
	testCases := []struct {
		name    string
		args    []string
		wantMap map[string]string
		wantOut string
		wantErr error
	}{
		{
			name: "set a single alias",
			args: []string{
				"answer-to-life=\"echo 43\"",
			},
			wantMap: map[string]string{
				"answer-to-life": "echo 43",
			},
		},
		{
			name: "set multiple aliases",
			args: []string{
				"answer-to-life=\"echo 43\"",
				"ls=\"ls -l\"",
			},
			wantMap: map[string]string{
				"answer-to-life": "echo 42",
				"ls":             "ls -l",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			builtins.AliasMap = make(map[string]string)

			var out bytes.Buffer

			err := builtins.Alias(&out, tc.args...)
			if err != tc.wantErr {
				t.Fatalf("expected %v, got %v error", tc.wantErr, err)
			}
		})
	}
}
