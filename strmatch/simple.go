package strmatch

import (
	"strings"
)

func simpleMatch(in string, slice []string, opt option) (res []Matched) {
	if !opt.caseSensitive {
		in = strings.ToLower(in)
		for i := range slice {
			slice[i] = strings.ToLower(slice[i])
		}
	}

	var idx int
	rIn := []rune(in)
	for i := range slice {
		s := []rune(slice[i])
	M:
		for j := range s {
			for k := range rIn[idx:] {
				if s[j] == rIn[k] {
					idx++
					if idx == len(rIn) {
						// log.Printf("%s matched", slice[i])
						break M
					}
					break
				}
			}
		}
	}
	return
}

func simpleMatch2(in string, slice []string, opt option) (res []Matched) {
	if !opt.caseSensitive {
		in = strings.ToLower(in)
		for i := range slice {
			slice[i] = strings.ToLower(slice[i])
		}
	}

	var idx int
	rIn := []rune(in)
	for i := range slice {
		s := []rune(slice[i])
	M:
		for j := range s {
			for k := range rIn[idx:] {
				if s[j] == rIn[k] {
					idx++
					if idx == len(rIn) {
						// log.Printf("%s matched", slice[i])
						break M
					}
					break
				}
			}
		}
	}
	return
}
