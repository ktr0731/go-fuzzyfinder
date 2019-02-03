package matching

import (
	"strings"

	"github.com/ktr0731/go-fuzzyfinder/scoring"
)

// Matched represents a result of FindAll.
type Matched struct {
	// Idx is the index of an item of the original slice which was used to
	// search matched strings.
	Idx int
	// Pos is the range of matched position.
	// [2]int represents a closed interval of a position.
	Pos [2]int
	// Score is the value that indicates how it similar to the input string.
	// The bigger Score, the more similar it is.
	Score int
}

type opt optFunc
type optFunc func(*option)

// option represents available options and its default values.
type option struct {
	caseSensitive bool
}

// WithCaseSensitive enables a case sensitive searching.
func WithCaseSensitive() opt {
	return func(o *option) {
		o.caseSensitive = true
	}
}

// FindAll tries to find out sub-strings from slice that match the passed argument in.
func FindAll(in string, slice []string, opts ...opt) []Matched {
	var opt option
	for _, o := range opts {
		o(&opt)
	}
	return match(in, slice, opt)
}

// match iterates each string of slice for check whether it is matched to the input string.
func match(input string, slice []string, opt option) (res []Matched) {
	if !opt.caseSensitive {
		input = strings.ToLower(input)
	}
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
						// TODO: 引数と順番をあわせる
						Score: scoring.Calculate(s, input),
					})
					break LINE_MATCHING
				}
			}
		}
	}
	return
}
