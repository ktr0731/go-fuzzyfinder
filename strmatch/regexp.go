package strmatch

import (
	"regexp"
)

func regexpMatch(in string, slice []string, opt option) (res []Matched) {
	pat := "("
	if !opt.caseSensitive {
		pat = "(?i)" + pat
	}
	for _, r := range []rune(in) {
		pat += regexp.QuoteMeta(string(r)) + ".*?"
	}
	pat += ")"
	re := regexp.MustCompile(pat)

	for i, s := range slice {
		pos := re.FindAllStringSubmatchIndex(s, -1)
		if pos == nil || len(pos[0]) != 4 {
			continue
		}
		res = append(res, Matched{
			Idx: i,
			Pos: [][2]int{{pos[0][2], pos[0][3]}},
		})
	}
	return
}
