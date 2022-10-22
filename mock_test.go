package fuzzyfinder

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nsf/termbox-go"
)

func Test_parseAttr(t *testing.T) {
	cases := map[string]struct {
		attr      termbox.Attribute
		isFg      bool
		expected  string
		willPanic bool
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
		"ColorGreen with reverse": {
			attr:     termbox.ColorGreen | termbox.AttrReverse,
			expected: "\x1b[7;48;5;2m",
		},
		"invalid color": {
			attr:      termbox.ColorWhite + 1,
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
			actual := parseAttr(c.attr, c.isFg)
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Errorf("diff found: \n%s\nexpected = %x, actual = %x", diff, c.expected, actual)
			}
		})
	}
}

// func Test_parseAttrV2(t *testing.T) {
// 	cases := map[string]struct {
// 		attr      tcell.AttrMask
// 		fg        tcell.Color
// 		bg        tcell.Color
// 		isBg      bool
// 		expected  string
// 		willPanic bool
// 	}{
// 		"ColorDefault": {
// 			fg:       tcell.ColorDefault,
// 			expected: "\x1b[39m",
// 		},
// 		"ColorDefault bg": {
// 			bg:       tcell.ColorDefault,
// 			isBg:     true,
// 			expected: "\x1b[49m",
// 		},
// 		"ColorGreen": {
// 			fg:       tcell.ColorGreen,
// 			expected: "\x1b[38;5;2m",
// 		},
// 		"ColorGreen with bold": {
// 			attr:     tcell.AttrBold,
// 			fg:       tcell.ColorGreen,
// 			expected: "\x1b[1;38;5;2m",
// 		},
// 		"ColorGreen with bold and underline": {
// 			attr:     tcell.AttrBold | tcell.AttrUnderline,
// 			fg:       tcell.ColorGreen,
// 			expected: "\x1b[1;4;38;5;2m",
// 		},
// 		"ColorGreen with reverse": {
// 			attr:     tcell.AttrReverse,
// 			fg:       tcell.ColorGreen,
// 			expected: "\x1b[7;38;5;2m",
// 		},
// 		"invalid color": {
// 			attr:      tcell.AttrInvalid,
// 			willPanic: true,
// 		},
// 	}
//
// 	for name, c := range cases {
// 		c := c
// 		t.Run(name, func(t *testing.T) {
// 			if c.willPanic {
// 				defer func() {
// 					if err := recover(); err == nil {
// 						t.Errorf("must panic")
// 					}
// 				}()
// 			}
// 			var actual string
// 			if c.isBg {
// 				actual = parseAttrV2(nil, &c.bg, c.attr)
// 			} else {
// 				actual = parseAttrV2(&c.fg, nil, c.attr)
// 			}
// 			if diff := cmp.Diff(c.expected, actual); diff != "" {
// 				t.Errorf("diff found: \n%s\nexpected = %x, actual = %x", diff, c.expected, actual)
// 			}
// 		})
// 	}
// }
