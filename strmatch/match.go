package strmatch

import (
	"strings"
)

func match(input string, _slice []string, opt option) (res []Matched) {
	slice := _slice
	if !opt.caseSensitive {
		s := make([]string, len(slice))
		input = strings.ToLower(input)
		for i := range _slice {
			s[i] = strings.ToLower(_slice[i])
		}
		slice = s
	}

	in := []rune(input)
	for idxOfSlice, s := range slice {
		var from int
		var idx int
	LINE_MATCHING:
		for i, r := range []rune(s) {
			if r == in[idx] {
				if idx == 0 {
					from = i
				}
				idx++
				if idx == len(in) {
					res = append(res, Matched{
						Idx: idxOfSlice,
						Pos: [2]int{from, i + 1},
					})
					break LINE_MATCHING
				}
			}
		}
	}
	return
}
