package builtins_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/raian621/csce-4600-project2/builtins"
)

func Test_Alias(t *testing.T) {
	testCases := []struct {
		name    string
		args    []string
		currMap map[string]string
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
				"answer-to-life=\"echo 42\"",
				"ls=\"ls -l\"",
			},
			wantMap: map[string]string{
				"answer-to-life": "echo 42",
				"ls":             "ls -l",
			},
		},
		{
			name: "alias with no aliases defined",
		},
		{
			name: "alias -p with no aliases defined",
			args: []string{"-p"},
		},
		{
			name: "alias -p with no aliases defined and alias assignments",
			args: []string{"-p", "a=alias"},
			wantMap: map[string]string{
				"a": "alias",
			},
		},
		{
			name: "print all aliases implicitly with no args",
			currMap: map[string]string{
				"a": "alias",
				"c": "cat",
				"l": "ls -l",
			},
			wantMap: map[string]string{
				"a": "alias",
				"c": "cat",
				"l": "ls -l",
			},
			wantOut: "alias a='alias'\nalias c='cat'\nalias l='ls -l'\n",
		},
		{
			name: "print all aliases explicitly with flag -p",
			args: []string{"-p"},
			currMap: map[string]string{
				"a": "alias",
				"c": "cat",
				"l": "ls -l",
			},
			wantMap: map[string]string{
				"a": "alias",
				"c": "cat",
				"l": "ls -l",
			},
			wantOut: "alias a='alias'\nalias c='cat'\nalias l='ls -l'\n",
		},
		{
			name: "print specified aliases with one nonexistent alias",
			args: []string{"a", "b", "c"},
			currMap: map[string]string{
				"a": "alias",
				"c": "cat",
				"l": "ls -l",
			},
			wantMap: map[string]string{
				"a": "alias",
				"c": "cat",
				"l": "ls -l",
			},
			wantOut: "alias a='alias'\nalias 'b' not found\nalias c='cat'\n",
			wantErr: builtins.ErrUndefinedAlias,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			builtins.AliasMap = tc.currMap
			if builtins.AliasMap == nil {
				builtins.AliasMap = make(map[string]string)
			}

			var out bytes.Buffer

			err := builtins.Alias(&out, tc.args...)
			if err != tc.wantErr {
				t.Fatalf("expected '%v', got '%v' error", tc.wantErr, err)
			}

			if out.String() != tc.wantOut {
				t.Fatalf("expected %s, got %s out", tc.wantOut, out.String())
			}

			if !reflect.DeepEqual(builtins.AliasMap, tc.wantMap) && (len(builtins.AliasMap) != 0 || len(tc.wantMap) != 0) {
				t.Fatalf("expected %v, got %v aliases", tc.wantMap, builtins.AliasMap)
			}
		})
	}

}
