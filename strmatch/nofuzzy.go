package strmatch

import "strings"

func noFuzzy(in string, slice []string, opt option) (res []Matched) {
	if !opt.caseSensitive {
		in = strings.ToLower(in)
	}
	for i, s := range slice {
		if !opt.caseSensitive {
			s = strings.ToLower(s)
		}
		pos := strings.Index(s, in)
		if pos == -1 {
			continue
		}
		res = append(res, Matched{
			Idx: i,
			Pos: [][2]int{{pos, pos + len(in)}},
		})
	}
	return
}
