package fuzzyfinder

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	runewidth "github.com/mattn/go-runewidth"
)

type simScreen tcell.SimulationScreen

// TerminalMock is a mocked terminal for testing.
// Most users should use it by calling UseMockedTerminal.
type TerminalMock struct {
	simScreen
	resultMu sync.RWMutex
	result   string
}

// UseMockedTerminal switches the terminal, which is used from
// this package to a mocked one.
func UseMockedTerminal() *TerminalMock {
	return defaultFinder.UseMockedTerminal()
}

func (t *TerminalMock) GetResult() string {
	var s string

	// set cursor for snapshot test
	setCursor := func() {
		cursorX, cursorY, _ := t.GetCursor()
		mainc, _, _, _ := t.GetContent(cursorX, cursorY)
		if mainc == ' ' {
			t.SetContent(cursorX, cursorY, '\u2588', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorDefault))
		} else {
			t.SetContent(cursorX, cursorY, mainc, nil, tcell.StyleDefault.Background(tcell.ColorWhite))
		}
		t.Show()
	}

	setCursor()

	t.resultMu.Lock()

	cells, width, height := t.GetContents()

	for h := 0; h < height; h++ {
		prevFg, prevBg := tcell.ColorDefault, tcell.ColorDefault
		for w := 0; w < width; w++ {
			cell := cells[h*width+w]
			fg, bg, attr := cell.Style.Decompose()
			var fgReset bool
			if fg != prevFg {
				s += "\x1b\x5b\x6d" // Reset previous color.
				s += parseAttr(&fg, nil, attr)
				prevFg = fg
				prevBg = tcell.ColorDefault
				fgReset = true
			}
			if bg != prevBg {
				if !fgReset {
					s += "\x1b\x5b\x6d" // Reset previous color.
					prevFg = tcell.ColorDefault
				}
				s += parseAttr(nil, &bg, attr)
				prevBg = bg
			}
			s += string(cell.Runes[:])
			rw := runewidth.RuneWidth(cell.Runes[0])
			if rw != 0 {
				w += rw - 1
			}
		}
		s += "\n"
	}
	s += "\x1b\x5b\x6d" // Reset previous color.

	t.resultMu.Unlock()

	return s
}

func (t *TerminalMock) SetEvents(events ...tcell.Event) {
	for _, event := range events {
		switch event.(type) {
		case *tcell.EventKey:
			ek := event.(*tcell.EventKey)
			t.InjectKey(ek.Key(), ek.Rune(), ek.Modifiers())
		case *tcell.EventResize:
			er := event.(*tcell.EventResize)
			w, h := er.Size()
			t.SetSize(w, h)
		}
	}
}

func (f *finder) UseMockedTerminal() *TerminalMock {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		panic(err)
	}
	m := &TerminalMock{
		simScreen: screen,
	}
	f.term = m
	return m
}

// as an escape sequence.
func parseAttr(fg, bg *tcell.Color, attr tcell.AttrMask) string {
	if attr == tcell.AttrInvalid {
		panic("invalid attribute")
	}

	var buf bytes.Buffer

	buf.WriteString("\x1b[")
	parseAttrMask := func() {
		if attr >= tcell.AttrUnderline {
			buf.WriteString("4;")
			attr -= tcell.AttrUnderline
		}
		if attr >= tcell.AttrReverse {
			buf.WriteString("7;")
			attr -= tcell.AttrReverse
		}
		if attr >= tcell.AttrBold {
			buf.WriteString("1;")
			attr -= tcell.AttrBold
		}
	}

	if fg != nil || bg != nil {
		isFg := fg != nil && bg == nil

		if isFg {
			parseAttrMask()
			if *fg == tcell.ColorDefault {
				buf.WriteString("39")
			} else {
				fmt.Fprintf(&buf, "38;5;%d", toAnsi3bit(*fg))
			}
		} else {
			if *bg == tcell.ColorDefault {
				buf.WriteString("49")
			} else {
				fmt.Fprintf(&buf, "48;5;%d", toAnsi3bit(*bg))
			}
		}
		buf.WriteString("m")
	}
	return buf.String()
}

func toAnsi3bit(color tcell.Color) int {
	colors := []tcell.Color{
		tcell.ColorBlack, tcell.ColorRed, tcell.ColorGreen, tcell.ColorYellow, tcell.ColorBlue, tcell.ColorDarkMagenta, tcell.ColorDarkCyan, tcell.ColorWhite,
	}
	for i, c := range colors {
		if c == color {
			return i
		}
	}
	return 0
}
