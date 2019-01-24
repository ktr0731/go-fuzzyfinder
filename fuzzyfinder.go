package fuzzyfinder

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/ktr0731/go-fuzzyfinder/strmatch"
	runewidth "github.com/mattn/go-runewidth"
	termbox "github.com/nsf/termbox-go"
	"github.com/pkg/errors"
)

var (
	// ErrAbort may be returned from FInd* functions if there are no selections.
	ErrAbort   = errors.New("abort")
	errEntered = errors.New("entered")
)

var defaultFinder = &finder{}

type state struct {
	items      []string           // All item names.
	allMatched []strmatch.Matched // All items.
	matched    []strmatch.Matched

	// x is the current index of the input line.
	x int
	// cursorX is the position of input line.
	// Note that cursorX is the actual width of input runes.
	cursorX int

	// The current index of filtered items (matched).
	// The initial state is 0.
	y int
	// cursorY is the position of item line.
	// Note that the max size of cursorY depends on max height.
	cursorY int

	input []rune

	// selections holds whether a key is selected or not. Each key is
	// an index of a item (Matched.Idx). Each value represents the item is
	// selected if true.
	selection map[int]bool
}

type finder struct {
	stateMu   sync.Mutex
	state     state
	drawTimer *time.Timer
	opt       *opt
}

func (f *finder) initFinder(items []string, matched []strmatch.Matched, opts []option) error {
	*f = finder{}
	if err := term.init(); err != nil {
		return errors.Wrap(err, "failed to initialize termbox")
	}

	var opt opt
	for _, o := range opts {
		o(&opt)
	}
	f.opt = &opt

	if opt.multi {
		f.state.selection = map[int]bool{}
	}

	f.state.items = items
	f.state.matched = matched
	f.state.allMatched = matched
	f.drawTimer = time.AfterFunc(0, func() {
		f._draw()
		f._drawPreview()
		term.flush()
	})
	f.drawTimer.Stop()
	return nil
}

// _draw is used from draw with a timer.
func (f *finder) _draw() {
	width, height := term.size()
	term.clear(termbox.ColorDefault, termbox.ColorDefault)

	maxWidth := width
	if f.opt.previewFunc != nil {
		maxWidth = width/2 - 1
	}

	// input line
	term.setCell(0, height-1, '>', termbox.ColorBlue, termbox.ColorDefault)
	var r rune
	var w int
	for _, r = range f.state.input {
		// Add a space between '>' and runes.
		term.setCell(2+w, height-1, r, termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
		w += runewidth.RuneWidth(r)
	}
	term.setCursor(2+f.state.cursorX, height-1)

	// Number line
	for i, r := range fmt.Sprintf("%d/%d", len(f.state.matched), len(f.state.items)) {
		term.setCell(2+i, height-2, r, termbox.ColorYellow, termbox.ColorDefault)
	}

	// Item lines
	itemAreaHeight := height - 2 - 1
	matched := f.state.matched
	offset := f.state.cursorY
	y := f.state.y
	// From the first (the most bottom) item in the item lines to the end.
	matched = matched[y-offset:]

	for i, m := range matched {
		if i > itemAreaHeight {
			break
		}
		if i == f.state.cursorY {
			term.setCell(0, height-3-i, '>', termbox.ColorRed, termbox.ColorBlack)
			term.setCell(1, height-3-i, ' ', termbox.ColorRed, termbox.ColorBlack)
		}

		if f.opt.multi {
			if f.state.selection[m.Idx] {
				term.setCell(1, height-3-i, '>', termbox.ColorRed, termbox.ColorBlack)
			}
		}

		var posIdx int
		w := 2
		for j, r := range f.state.items[m.Idx] {
			fg := termbox.ColorDefault
			bg := termbox.ColorDefault
			// Highlight selected strings.
			if len(m.Pos) != 0 {
				from, to := m.Pos[0][0], m.Pos[0][1]
				if !(from == 0 && to == 0) && (j >= from && j < to) {
					if unicode.ToLower(f.state.input[posIdx]) == unicode.ToLower(r) {
						fg |= termbox.ColorGreen
						posIdx++
					}
				}
			}
			if i == f.state.cursorY {
				fg |= termbox.AttrBold | termbox.ColorYellow
				bg = termbox.ColorBlack
			}

			rw := runewidth.RuneWidth(r)
			// Shorten item cells.
			if w+rw+2 > maxWidth {
				term.setCell(w, height-3-i, '.', fg, bg)
				term.setCell(w+1, height-3-i, '.', fg, bg)
				w += 2
				break
			} else {
				term.setCell(w, height-3-i, r, fg, bg)
				w += rw
			}
		}
	}
}

func (f *finder) _drawPreview() {
	if f.opt.previewFunc == nil {
		return
	}

	width, height := term.size()
	var idx int
	if len(f.state.matched) == 0 {
		idx = -1
	} else {
		idx = f.state.matched[f.state.y].Idx
	}

	sp := strings.Split(f.opt.previewFunc(idx, width, height), "\n")
	prevLines := make([][]rune, 0, len(sp))
	for _, s := range sp {
		prevLines = append(prevLines, []rune(s))
	}

	// top line
	for i := width / 2; i < width; i++ {
		var r rune
		if i == width/2 {
			r = '┌'
		} else if i == width-1 {
			r = '┐'
		} else {
			r = '─'
		}
		term.setCell(i, 0, r, termbox.ColorBlack, termbox.ColorDefault)
	}
	// bottom line
	for i := width / 2; i < width; i++ {
		var r rune
		if i == width/2 {
			r = '└'
		} else if i == width-1 {
			r = '┘'
		} else {
			r = '─'
		}
		term.setCell(i, height-1, r, termbox.ColorBlack, termbox.ColorDefault)
	}
	// Start with h=1 to exclude each corner rune.
	const vline = '│'
	var wvline = runewidth.RuneWidth(vline)
	for h := 1; h < height-1; h++ {
		w := width / 2
		for i := width / 2; i < width; i++ {
			switch {
			// Left vertical line.
			case i == width/2:
				term.setCell(i, h, vline, termbox.ColorBlack, termbox.ColorDefault)
				w += wvline
			// Right vertical line.
			case i == width-1:
				term.setCell(i, h, vline, termbox.ColorBlack, termbox.ColorDefault)
				w += wvline
			// Spaces between left and right vertical lines.
			case w == width/2+wvline, w == width-1-wvline:
				term.setCell(w, h, ' ', termbox.ColorDefault, termbox.ColorDefault)
				w++
			default: // Preview text
				if h-1 >= len(prevLines) {
					w++
					continue
				}
				j := i - width/2 - 2 // Two spaces.
				l := prevLines[h-1]
				if j >= len(l) {
					w++
					continue
				}
				rw := runewidth.RuneWidth(l[j])
				if w+rw > width-1-2 {
					term.setCell(w, h, '.', termbox.ColorDefault, termbox.ColorDefault)
					term.setCell(w+1, h, '.', termbox.ColorDefault, termbox.ColorDefault)
					w += 2
					continue
				}

				term.setCell(w, h, l[j], termbox.ColorDefault, termbox.ColorDefault)
				w += rw
			}
		}
	}
}

func (f *finder) draw(d time.Duration) {
	f.stateMu.Lock()
	defer f.stateMu.Unlock()

	if isInTesting() {
		// Don't use goroutine scheduling.
		f._draw()
		f._drawPreview()
		term.flush()
	} else {
		f.drawTimer.Reset(d)
	}
}

func (f *finder) readKey() error {
	switch e := term.pollEvent(); e.Type {
	case termbox.EventKey:
		switch e.Key {
		case termbox.KeyEsc, termbox.KeyCtrlC, termbox.KeyCtrlD:
			return ErrAbort
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			if len(f.state.input) == 0 {
				return nil
			}
			if f.state.x == 0 {
				return nil
			}
			// Remove the latest input rune.
			f.state.cursorX -= runewidth.RuneWidth(f.state.input[len(f.state.input)-1])
			f.state.x--
			f.state.input = f.state.input[0 : len(f.state.input)-1]
		case termbox.KeyDelete:
			if f.state.x == len(f.state.input) {
				return nil
			}
			x := f.state.x
			f.state.input = append(f.state.input[:x], f.state.input[x+1:]...)
		case termbox.KeyEnter:
			return errEntered
		case termbox.KeyArrowLeft, termbox.KeyCtrlB:
			if f.state.x > 0 {
				f.state.cursorX -= runewidth.RuneWidth(f.state.input[f.state.x-1])
				f.state.x--
			}
		case termbox.KeyArrowRight, termbox.KeyCtrlF:
			if f.state.x < len(f.state.input) {
				f.state.cursorX += runewidth.RuneWidth(f.state.input[f.state.x])
				f.state.x++
			}
		case termbox.KeyCtrlA:
			f.state.cursorX = 0
			f.state.x = 0
		case termbox.KeyCtrlE:
			f.state.cursorX = runewidth.StringWidth(string(f.state.input))
			f.state.x = len(f.state.input)
		case termbox.KeyCtrlW:
			in := f.state.input[:f.state.x]
			inStr := string(in)
			pos := strings.LastIndex(strings.TrimRightFunc(inStr, unicode.IsSpace), " ")
			if pos == -1 {
				f.state.input = []rune{}
				f.state.cursorX = 0
				f.state.x = 0
				return nil
			}
			pos = utf8.RuneCountInString(inStr[:pos])
			newIn := f.state.input[:pos+1]
			f.state.input = newIn
			f.state.cursorX = runewidth.StringWidth(string(newIn))
			f.state.x = len(newIn)
		case termbox.KeyCtrlU:
			f.state.input = f.state.input[f.state.x:]
			f.state.cursorX = 0
			f.state.x = 0
		case termbox.KeyArrowUp, termbox.KeyCtrlK, termbox.KeyCtrlP:
			f.stateMu.Lock()
			defer f.stateMu.Unlock()

			if f.state.y+1 < len(f.state.matched) {
				f.state.y++
			}
			_, height := term.size()
			if f.state.cursorY+1 < height-2 && f.state.cursorY+1 < len(f.state.matched) {
				f.state.cursorY++
			}
		case termbox.KeyArrowDown, termbox.KeyCtrlJ, termbox.KeyCtrlN:
			if f.state.y > 0 {
				f.state.y--
			}
			if f.state.cursorY-1 >= 0 {
				f.state.cursorY--
			}
		case termbox.KeyTab:
			if !f.opt.multi {
				return nil
			}
			idx := f.state.matched[f.state.y].Idx
			if s, ok := f.state.selection[idx]; ok {
				f.state.selection[idx] = !s
			} else {
				f.state.selection[idx] = true
			}
			if f.state.y > 0 {
				f.state.y--
			}
			if f.state.cursorY > 0 {
				f.state.cursorY--
			}
		default:
			if e.Key == termbox.KeySpace {
				e.Ch = ' '
			}
			if e.Ch != 0 {
				x := f.state.x
				f.state.input = append(f.state.input[:x], append([]rune{e.Ch}, f.state.input[x:]...)...)
				f.state.cursorX += runewidth.RuneWidth(e.Ch)
				f.state.x++
			}
		}
	case termbox.EventResize:
		f.draw(200 * time.Millisecond)
	}
	return nil
}

func (f *finder) filter() {
	f.stateMu.Lock()
	defer f.stateMu.Unlock()

	if len(f.state.input) == 0 {
		f.state.matched = f.state.allMatched
		return
	}

	// TODO: If input is not delete operation, it is able to
	// reduce total iteration.
	matchedItems := strmatch.FindAll(string(f.state.input), f.state.items, strmatch.WithAlgo(strmatch.AlgoRegExp))
	var prevSelectedItemIdx int
	if len(f.state.matched) != 0 {
		prevSelectedItemIdx = f.state.matched[f.state.y].Idx
	}
	f.state.matched = matchedItems
	if len(f.state.matched) == 0 {
		f.state.cursorY = 0
		f.state.y = 0
		return
	}

	switch {
	case f.state.cursorY >= len(f.state.matched):
		f.state.cursorY = len(f.state.matched) - 1
		f.state.y = len(f.state.matched) - 1
	case f.state.y >= len(f.state.matched):
		f.state.y = len(f.state.matched) - 1
	// Update the current index and cursor to the selected track.
	default:
		_, height := term.size()
		itemAreaHeight := height - 2 - 1
		for i, m := range f.state.matched {
			if m.Idx == prevSelectedItemIdx {
				f.state.y = i
				if i > itemAreaHeight {
					f.state.cursorY = itemAreaHeight
				} else if i < 0 {
					f.state.cursorY = 0
				} else {
					f.state.cursorY = i
				}
				break
			}
		}
	}
}

func (f *finder) find(items []string, matched []strmatch.Matched, opts []option) ([]int, error) {
	if err := f.initFinder(items, matched, opts); err != nil {
		return nil, errors.Wrap(err, "failed to initialize the fuzzy finder")
	}
	defer term.close()

	for {
		f.draw(0)
		err := f.readKey()
		switch {
		case err == ErrAbort:
			return nil, ErrAbort
		case err == errEntered:
			f.stateMu.Lock()
			defer f.stateMu.Unlock()

			if len(f.state.matched) == 0 {
				return nil, ErrAbort
			}
			if f.opt.multi {
				if len(f.state.selection) == 0 {
					return []int{f.state.matched[f.state.y].Idx}, nil
				}
				idxs := make([]int, 0, len(f.state.selection))
				for idx, selected := range f.state.selection {
					if selected {
						idxs = append(idxs, idx)
					}
				}
				return idxs, nil
			}
			return []int{f.state.matched[f.state.y].Idx}, nil
		case err != nil:
			return nil, errors.Wrap(err, "failed to read a key")
		}
		f.filter()
	}
}

func isInTesting() bool {
	return flag.Lookup("test.v") != nil
}

// Find displays a UI that provides fuzzy finding against to the passed slice.
// The argument slice must be a slice type. If it is not a slice, Find returns
// an error. itemFunc is called by the length of slice. previewFunc is called
// when the cursor which points the current selected item is changed.
// If itemFunc is nil, Find returns a error. If previewFunc is nil, preview
// feature is disabled.
//
// itemFunc receives an argument i. It is the index of the item currently
// selected.
//
// Find returns ErrAbort if a call of Find is finished with no selection.
func Find(slice interface{}, itemFunc func(i int) string, opts ...option) (int, error) {
	res, err := find(slice, itemFunc, opts)
	if err != nil {
		return 0, err
	}
	return res[0], err
}

// FindMulti is nearly same as the Find. The only one difference point from
// Find is the user can select multiple items at once by tab key.
func FindMulti(slice interface{}, itemFunc func(i int) string, opts ...option) ([]int, error) {
	return find(slice, itemFunc, append(opts, withMulti()))
}

func find(slice interface{}, itemFunc func(i int) string, opts []option) ([]int, error) {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return nil, errors.Errorf("the first argument must be a slice, but got %T", slice)
	}
	if itemFunc == nil {
		return nil, errors.New("itemFunc must not be nil")
	}

	sliceLen := rv.Len()
	items := make([]string, sliceLen)
	matchedItems := make([]strmatch.Matched, sliceLen)
	for i := 0; i < sliceLen; i++ {
		items[i] = itemFunc(i)
		matchedItems[i] = strmatch.Matched{Idx: i}
	}

	return defaultFinder.find(items, matchedItems, opts)
}
