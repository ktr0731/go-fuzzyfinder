package strmatch

import "strings"

// match iterates each string of slice for check whether it is matched to the input string.
func match(input string, slice []string, opt option) (res []Matched) {
	in := []rune(input)
	for idxOfSlice, s := range slice {
		var from, idx int
		if !opt.caseSensitive {
			s = strings.ToLower(s)
		}
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
