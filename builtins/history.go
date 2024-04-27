package builtins

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var (
	HistoryList            []string = make([]string, 0)
	MaxHistory             int      = 420
	MinIndex               int      = 1
	ErrInvalidHistoryRange error    = errors.New("invalid range of history entries")
)

// https://www.gnu.org/software/bash/manual/html_node/Bash-History-Builtins.html#index-history
// I'm not going to implement any of the disk-persistent history file slop though.
/*
FLAGS:
  -c : Clears the history list.
  -d : Deletes an entry / entries from the history list.
	     To delete an entry, simply supply a valid entry after the -d flag.
			 To delete a range of entries, enter the start and end entry numbers,
			 seperated by a dash (ex. history -d 69-100). The specified range of
			 entries will be deleted from start to end inclusive.

NOTES:
  Bash's history command (at least on my machine) only stores 1000 entries. Additionally,
	the entries maintain their indices; Whenever whenever the 1001st entry is added,
	the 1st entry is removed and the 2nd entry becomes the first entry, but keeps its
	index of 2. This process is repeated for the 1000 + nth entry as well. Also,
	if a command is the same as the command (character for character) in the last
	entry of the history, it is not added to the history.

	If an entry is deleted normally though, the indices of the list are not
	maintained, every index after the deleted index or indices will be decremented
	by 1 or by however many entries were deleted.
*/
// may return strconv.ErrSyntax errors from attempts to parse ints
func History(w io.Writer, args ...string) error {
	if len(args) > 2 {
		fmt.Fprintln(w, "usage: history [-d start[-end]|-c]")
		return ErrInvalidArgCount
	}

	if len(args) == 2 && args[0] == "-d" {
		if i := strings.IndexByte(args[1], '-'); i != -1 {
			// handle entry range inputs
			startStr, endStr := args[1][:i], args[1][i+1:]
			start, err := strconv.ParseInt(startStr, 10, 32)
			if err != nil {
				fmt.Fprintln(w, "invalid start entry")
				return err
			}
			end, err := strconv.ParseInt(endStr, 10, 32)
			if err != nil {
				fmt.Fprintln(w, "invalid end entry")
				return err
			}

			lastIndex := int64(len(HistoryList) + MinIndex)
			if start < int64(MinIndex) || start >= lastIndex || start >= end {
				fmt.Fprintln(w, "invalid entry range")
				return ErrInvalidHistoryRange
			}
			if end >= lastIndex {
				fmt.Fprintln(w, "invalid entry range")
				return ErrInvalidHistoryRange
			}

			start -= int64(MinIndex)
			end -= int64(MinIndex)

			// remove range [start, end] (inclusive)
			HistoryList = append(HistoryList[:start], HistoryList[end+1:]...)
		} else {
			// handle single entry inputs
			index, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				fmt.Fprintln(w, "invalid entry")
				return err
			}

			lastIndex := int64(len(HistoryList) + MinIndex)
			if index < int64(MinIndex) || index >= lastIndex {
				fmt.Fprintln(w, "invalid entry range")
				return ErrInvalidHistoryRange
			}

			index -= int64(MinIndex)

			HistoryList = append(HistoryList[:index], HistoryList[index+1:]...)
		}
	} else if len(args) == 1 && args[0] == "-c" {
		// bruh
		HistoryList = make([]string, 0)
	} else if len(args) == 0 {
		// print out history
		for i, entry := range HistoryList {
			fmt.Fprintf(w, "%5d  %s\n", i+MinIndex, entry)
		}
	} else {
		fmt.Fprintln(w, "usage: history [-d start[-end]|-c]")
		return ErrUnknownArg
	}

	return nil
}

func AddHistoryEntry(entry string) {
	if len(HistoryList) != 0 && entry == HistoryList[len(HistoryList)-1] {
		return
	}

	if len(HistoryList) == MaxHistory {
		HistoryList = HistoryList[1:] // remove first entry in the list to make room
		MinIndex += 1
	}

	HistoryList = append(HistoryList, entry)
}
