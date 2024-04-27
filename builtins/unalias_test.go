package builtins_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/raian621/csce-4600-project2/builtins"
)

func TestUnalias(t *testing.T) {
	testCases := []struct {
		name    string
		args    []string
		currMap map[string]string
		wantMap map[string]string
		wantOut string
		wantErr error
	}{
		{
			name:    "no arguments should return an error and print an error message",
			wantMap: make(map[string]string), // empty map
			wantOut: "usage: unalias [-a] name [name ...]\n",
			wantErr: builtins.ErrInvalidArgCount,
		},
		{
			name: "unalias EVERYTHING",
			args: []string{"-a"},
			currMap: map[string]string{
				"a": "alias",
				"c": "cat",
				"l": "ls -l",
			},
			wantMap: make(map[string]string), // emptied map
		},
		{
			name: "unalias some aliases",
			args: []string{"c"},
			currMap: map[string]string{
				"a": "alias",
				"c": "cat",
				"l": "ls -l",
			},
			wantMap: map[string]string{
				"a": "alias",
				"l": "ls -l",
			},
		},
		{
			name: "unalias undefined alias",
			args: []string{"d"},
			currMap: map[string]string{
				"a": "alias",
				"c": "cat",
				"l": "ls -l",
			},
			wantOut: "alias 'd' not found\n",
			wantMap: map[string]string{
				"a": "alias",
				"c": "cat",
				"l": "ls -l",
			},
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

			err := builtins.Unalias(&out, tc.args...)
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
