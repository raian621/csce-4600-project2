package builtins_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/raian621/csce-4600-project2/builtins"
)

func TestHistory(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		currList       []string
		currMinIndex   int
		currMaxHistory int
		wantOut        string
		wantList       []string
		wantMinIndex   int
		wantErr        error
	}{
		{
			name: "print out history",
			args: []string{},
			currList: []string{
				"echo hello",
				"echo world",
				"echo ligma",
			},
			currMinIndex:   1,
			currMaxHistory: 20,
			wantOut: `    1  echo hello
    2  echo world
    3  echo ligma
`,
			wantList: []string{
				"echo hello",
				"echo world",
				"echo ligma",
			},
			wantMinIndex: 1,
		},
		{
			name: "remove a single entry",
			args: []string{"-d", "2"},
			currList: []string{
				"echo hello",
				"echo world",
				"echo ligma",
			},
			currMinIndex:   1,
			currMaxHistory: 20,
			wantList: []string{
				"echo hello",
				"echo ligma",
			},
			wantMinIndex: 1,
		},
		{
			name: "try to remove an invalid entry",
			args: []string{"-d", "asdf"},
			currList: []string{
				"echo hello",
				"echo world",
				"echo ligma",
			},
			currMinIndex:   1,
			currMaxHistory: 20,
			wantList: []string{
				"echo hello",
				"echo world",
				"echo ligma",
			},
			wantOut:      "invalid entry\n",
			wantMinIndex: 1,
			wantErr:      errors.New("strconv.ParseInt: parsing \"asdf\": invalid syntax"),
		},
		{
			name: "try to remove a nonexistent entry",
			args: []string{"-d", "4"},
			currList: []string{
				"echo hello",
				"echo world",
				"echo ligma",
			},
			currMinIndex:   1,
			currMaxHistory: 20,
			wantList: []string{
				"echo hello",
				"echo world",
				"echo ligma",
			},
			wantOut:      "invalid entry range\n",
			wantMinIndex: 1,
			wantErr:      builtins.ErrInvalidHistoryRange,
		},
		{
			name: "remove a range of entries",
			args: []string{"-d", "2-3"},
			currList: []string{
				"echo hello",
				"echo world",
				"echo ligma",
			},
			currMinIndex:   1,
			currMaxHistory: 20,
			wantList: []string{
				"echo hello",
			},
			wantMinIndex: 1,
		},
		{
			name:           "too many args",
			args:           []string{"-d", "-c", "sdfhj"},
			currList:       []string{},
			currMinIndex:   1,
			currMaxHistory: 20,
			wantList:       []string{},
			wantOut:        "usage: history [-d start[-end]|-c]\n",
			wantMinIndex:   1,
		},
		{
			name:           "try to remove an invalid start entry",
			args:           []string{"-d", "yo-3"},
			currList:       []string{},
			currMinIndex:   1,
			currMaxHistory: 20,
			wantList:       []string{},
			wantOut:        "invalid start entry\n",
			wantMinIndex:   1,
			wantErr:        errors.New("strconv.ParseInt: parsing \"yo\": invalid syntax"),
		},
		{
			name:           "try to remove an invalid start entry",
			args:           []string{"-d", "1-asdf"},
			currList:       []string{},
			currMinIndex:   1,
			currMaxHistory: 20,
			wantList:       []string{},
			wantOut:        "invalid end entry\n",
			wantMinIndex:   1,
			wantErr:        errors.New("strconv.ParseInt: parsing \"asdf\": invalid syntax"),
		},
		{
			name:           "try to remove a start entry outside the valid range",
			args:           []string{"-d", "1-2"},
			currList:       []string{},
			currMinIndex:   1,
			currMaxHistory: 20,
			wantList:       []string{},
			wantOut:        "invalid entry range\n",
			wantMinIndex:   1,
			wantErr:        builtins.ErrInvalidHistoryRange,
		},
		{
			name:           "try to remove an end entry outside the valid range",
			args:           []string{"-d", "1-2"},
			currList:       []string{"echo hello"},
			currMinIndex:   1,
			currMaxHistory: 20,
			wantList:       []string{"echo hello"},
			wantOut:        "invalid entry range\n",
			wantMinIndex:   1,
			wantErr:        builtins.ErrInvalidHistoryRange,
		},
		{
			name:           "clear history list",
			args:           []string{"-c"},
			currList:       []string{"echo hello"},
			currMinIndex:   1,
			currMaxHistory: 20,
			wantList:       []string{},
			wantMinIndex:   1,
		},
		{
			name:           "invalid flag",
			args:           []string{"-g"},
			currList:       []string{},
			currMinIndex:   1,
			currMaxHistory: 20,
			wantOut:        "usage: history [-d start[-end]|-c]\n",
			wantList:       []string{},
			wantMinIndex:   1,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer

			builtins.HistoryList = tc.currList
			err := builtins.History(&out, tc.args...)

			if err != nil && tc.wantErr != nil && err.Error() != tc.wantErr.Error() {
				t.Fatalf("expected '%v', got '%v' error", tc.wantErr, err)
			}

			if out.String() != tc.wantOut {
				t.Fatalf("expected '%s', got '%s' output", tc.wantOut, out.String())
			}

			if (len(tc.wantList) > 0 || len(builtins.HistoryList) > 0) && !reflect.DeepEqual(tc.wantList, builtins.HistoryList) {
				t.Fatalf("expected '%v', got '%v' output", tc.wantList, builtins.HistoryList)
			}

			if builtins.MinIndex != tc.wantMinIndex {
				t.Fatalf("expected '%d', got '%d' min history index", tc.wantMinIndex, builtins.MinIndex)
			}
		})
	}
}

func TestAddHistoryEntry(t *testing.T) {
	testCases := []struct {
		name           string
		entry          string
		currList       []string
		currMinIndex   int
		currMaxHistory int
		wantMinIndex   int
		wantList       []string
	}{
		{
			name:           "simply add an entry",
			entry:          "echo hello",
			currList:       []string{},
			currMinIndex:   1,
			currMaxHistory: 10,
			wantMinIndex:   1,
			wantList:       []string{"echo hello"},
		},
		{
			name:  "add an entry when history list is full",
			entry: "echo CSCE-4600",
			currList: []string{
				"echo hello",
				"echo world",
				"echo ligma",
			},
			currMinIndex:   1,
			currMaxHistory: 3,
			wantMinIndex:   2,
			wantList: []string{
				"echo world",
				"echo ligma",
				"echo CSCE-4600",
			},
		},
		{
			name:  "add an entry that's the same as the previous entry",
			entry: "echo ligma",
			currList: []string{
				"echo hello",
				"echo world",
				"echo ligma",
			},
			currMinIndex:   1,
			currMaxHistory: 3,
			wantMinIndex:   1,
			wantList: []string{
				"echo hello",
				"echo world",
				"echo ligma",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			builtins.HistoryList = tc.currList
			builtins.MaxHistory = tc.currMaxHistory

			builtins.AddHistoryEntry(tc.entry)

			if (len(builtins.HistoryList) != 0 || len(tc.wantList) != 0) &&
				!reflect.DeepEqual(builtins.HistoryList, tc.wantList) {
				t.Fatalf("expected '%v', got '%v' output", tc.wantList, builtins.HistoryList)
			}
		})
	}
}
