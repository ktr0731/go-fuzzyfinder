package strmatch_test

import (
	"testing"

	"github.com/ktr0731/iv/fuzzyfinder/strmatch"
)

func TestAlgoNoFuzzy(t *testing.T) {
	cases := map[string]struct {
		in            string
		expected      string // If expected is empty, it means there are no matched strings.
		caseSensitive bool
	}{
		"normal":          {in: "ink", expected: "ink"},
		"case sensitive":  {in: "Twi", expected: "Twi", caseSensitive: true},
		"case sensitive2": {in: "twi", caseSensitive: true},
	}
	slice := []string{
		"WHITE ALBUM",
		"SOUND OF DESTINY",
		"Twinkle Snow",
	}
	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			var matched []strmatch.Matched
			if c.caseSensitive {
				matched = strmatch.FindAll(c.in, slice, strmatch.WithAlgo(strmatch.AlgoNoFuzzy), strmatch.WithCaseSensitive())
			} else {
				matched = strmatch.FindAll(c.in, slice, strmatch.WithAlgo(strmatch.AlgoNoFuzzy))
			}
			n := len(matched)
			if c.expected == "" {
				if n != 0 {
					t.Errorf("the result length must be 0, but got %d", n)
				}
				return
			}

			if n != 1 {
				t.Fatalf("the result length must be 1, but got %d", n)
			}
			m := matched[0]
			if m.Idx != 2 {
				t.Errorf("m.Idx must be equal to the index of 'Twinkle Snow', but got %d", m.Idx)
			}
			if n := len(m.Pos); n != 1 {
				t.Fatalf("AlgoNoFuzzy must return only one pos, but got %d", n)
			}
			from, to := m.Pos[0][0], m.Pos[0][1]
			if slice[2][from:to] != c.expected {
				t.Errorf("invalid range: from = %d, to = %d, content = %s, expected = %s", from, to, slice[2][from:to], c.expected)
			}
		})
	}
}

func BenchmarkNoFuzzy(b *testing.B) {
	in := "ink"
	slice := []string{
		"WHITE ALBUM",
		"SOUND OF DESTINY",
		"Twinkle Snow",
	}
	for i := 0; i < b.N; i++ {
		strmatch.FindAll(in, slice, strmatch.WithAlgo(strmatch.AlgoNoFuzzy))
	}
}
