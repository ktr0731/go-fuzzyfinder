package fuzzyfinder

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	termbox "github.com/nsf/termbox-go"
)

func Test_parseAttr(t *testing.T) {
	cases := map[string]struct {
		attr     termbox.Attribute
		isFg     bool
		expected string
	}{
		"ColorDefault": {
			attr:     termbox.ColorDefault,
			isFg:     true,
			expected: "\x1b[39m",
		},
		"ColorDefault bg": {
			attr:     termbox.ColorDefault,
			expected: "\x1b[49m",
		},
		"ColorGreen": {
			attr:     termbox.ColorGreen,
			expected: "\x1b[48;5;2m",
		},
		"ColorGreen with bold": {
			attr:     termbox.ColorGreen | termbox.AttrBold,
			expected: "\x1b[1;48;5;2m",
		},
		"ColorGreen with bold and underline": {
			attr:     termbox.ColorGreen | termbox.AttrBold | termbox.AttrUnderline,
			expected: "\x1b[4;1;48;5;2m",
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			actual := parseAttr(c.attr, c.isFg)
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Errorf("diff found: \n%s\nexpected = %x, actual = %x", diff, c.expected, actual)
			}
		})
	}
}
