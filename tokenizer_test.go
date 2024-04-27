package main

import (
	"reflect"
	"testing"
)

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
