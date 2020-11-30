package fuzzyfinder

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/google/go-cmp/cmp"
)

func Test_parseAttr(t *testing.T) {
	cases := map[string]struct {
		attr      tcell.AttrMask
		fg        tcell.Color
		bg        tcell.Color
		isBg      bool
		expected  string
		willPanic bool
	}{
		"ColorDefault": {
			fg:       tcell.ColorDefault,
			expected: "\x1b[39m",
		},
		"ColorDefault bg": {
			bg:       tcell.ColorDefault,
			isBg:     true,
			expected: "\x1b[49m",
		},
		"ColorGreen": {
			fg:       tcell.ColorGreen,
			expected: "\x1b[38;5;2m",
		},
		"ColorGreen with bold": {
			attr:     tcell.AttrBold,
			fg:       tcell.ColorGreen,
			expected: "\x1b[1;38;5;2m",
		},
		"ColorGreen with bold and underline": {
			attr:     tcell.AttrBold | tcell.AttrUnderline,
			fg:       tcell.ColorGreen,
			expected: "\x1b[4;1;38;5;2m",
		},
		"ColorGreen with reverse": {
			attr:     tcell.AttrReverse,
			fg:       tcell.ColorGreen,
			expected: "\x1b[7;38;5;2m",
		},
		"invalid color": {
			attr:      tcell.AttrInvalid,
			willPanic: true,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			if c.willPanic {
				defer func() {
					if err := recover(); err == nil {
						t.Errorf("must panic")
					}
				}()
			}
			var actual string
			if c.isBg {
				actual = parseAttr(nil, &c.bg, c.attr)
			} else {
				actual = parseAttr(&c.fg, nil, c.attr)
			}
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Errorf("diff found: \n%s\nexpected = %x, actual = %x", diff, c.expected, actual)
			}
		})
	}
}
