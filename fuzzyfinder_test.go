package fuzzyfinder_test

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gdamore/tcell/termbox"
	"github.com/google/go-cmp/cmp"
	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
)

var (
	update = flag.Bool("update", false, "update golden files")
	real   = flag.Bool("real", false, "display the actual layout to the terminal")
)

func init() {
	testing.Init()
	flag.Parse()
	if *update {
		os.RemoveAll(filepath.Join("testdata", "fixtures"))
		os.MkdirAll(filepath.Join("testdata", "fixtures"), 0755)
	}
}

func assertWithGolden(t *testing.T, f func(t *testing.T) string) {
	name := t.Name()
	r := strings.NewReplacer(
		"/", "-",
		" ", "_",
		"=", "-",
		"'", "",
		`"`, "",
		",", "",
	)
	normalizeFilename := func(name string) string {
		fname := r.Replace(strings.ToLower(name)) + ".golden"
		return filepath.Join("testdata", "fixtures", fname)
	}

	actual := f(t)

	fname := normalizeFilename(name)

	if *update {
		if err := ioutil.WriteFile(fname, []byte(actual), 0644); err != nil {
			t.Fatalf("failed to update the golden file: %s", err)
		}
		return
	}

	// Load the golden file.
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatalf("failed to load a golden file: %s", err)
	}
	expected := string(b)
	if runtime.GOOS == "windows" {
		expected = strings.Replace(expected, "\r\n", "\n", -1)
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("wrong result: \n%s", diff)
	}
}

type track struct {
	Name   string
	Artist string
	Album  string
}

var tracks = []*track{
	{"あの日自分が出て行ってやっつけた時のことをまだ覚えている人の為に", "", ""},
	{"ヒトリノ夜", "ポルノグラフィティ", "ロマンチスト・エゴイスト"},
	{"adrenaline!!!", "TrySail", "TAILWIND"},
	{"ソラニン", "ASIAN KUNG-FU GENERATION", "ソラニン"},
	{"closing", "AQUAPLUS", "WHITE ALBUM2"},
	{"glow", "keeno", "in the rain"},
	{"メーベル", "バルーン", "Corridor"},
	{"ICHIDAIJI", "ポルカドットスティングレイ", "一大事"},
	{"Catch the Moment", "LiSA", "Catch the Moment"},
}

func TestReal(t *testing.T) {
	if !*real {
		t.Skip("--real is disabled")
		return
	}
	_, err := fuzzyfinder.Find(
		tracks,
		func(i int) string {
			return tracks[i].Name
		},
		fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
			if i == -1 {
				return "not found"
			}
			return "Name: " + tracks[i].Name + "\nArtist: " + tracks[i].Artist
		}),
	)
	if err != nil {
		t.Fatalf("err is not nil: %s", err)
	}
}

func TestFind(t *testing.T) {
	cases := map[string]struct {
		events []termbox.Event
	}{
		"initial":          {},
		"input lo":         {runes("lo")},
		"input glow":       {runes("glow")},
		"arrow up-down":    {keys(termbox.KeyArrowUp, termbox.KeyArrowUp, termbox.KeyArrowDown)},
		"arrow left-right": {append(runes("ゆるふわ樹海ガール"), keys(termbox.KeyArrowLeft, termbox.KeyArrowLeft, termbox.KeyArrowLeft, termbox.KeyArrowRight)...)},
		"backspace":        {append(runes("adrenaline!!! -TV Ver.-"), keys(termbox.KeyBackspace, termbox.KeyBackspace)...)},
		"backspace empty":  {keys(termbox.KeyBackspace2, termbox.KeyBackspace2)},
		"backspace2":       {append(runes("オレンジ"), keys(termbox.KeyBackspace2, termbox.KeyBackspace2)...)},
		"delete":           {append(runes("オレンジ"), keys(termbox.KeyCtrlA, termbox.KeyDelete)...)},
		"delete empty":     {keys(termbox.KeyCtrlA, termbox.KeyDelete)},
		"ctrl-e":           {append(runes("恋をしたのは"), keys(termbox.KeyCtrlA, termbox.KeyCtrlE)...)},
		"ctrl-w":           {append(runes("ハロ / ハワユ"), key(termbox.KeyCtrlW))},
		"ctrl-w emtpy":     {keys(termbox.KeyCtrlW)},
		"ctrl-u":           {append(runes("恋をしたのは"), keys(termbox.KeyArrowLeft, termbox.KeyCtrlU, termbox.KeyArrowRight)...)},
		"long item":        {keys(termbox.KeyArrowUp, termbox.KeyArrowUp, termbox.KeyArrowUp)},
		"paging":           {keys(termbox.KeyArrowUp, termbox.KeyArrowUp, termbox.KeyArrowUp, termbox.KeyArrowUp, termbox.KeyArrowUp, termbox.KeyArrowUp, termbox.KeyArrowUp, termbox.KeyArrowUp)},
		"tab doesn't work": {keys(termbox.KeyTab)},
		"backspace doesnt change x if cursorX is 0": {append(runes("a"), keys(termbox.KeyCtrlA, termbox.KeyBackspace, termbox.KeyCtrlF)...)},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			c := c
			events := c.events

			f, term := fuzzyfinder.NewWithMockedTerminal()
			events = append(events, key(termbox.KeyEsc))
			term.SetEvents(events...)

			assertWithGolden(t, func(t *testing.T) string {
				_, err := f.Find(
					tracks,
					func(i int) string {
						return tracks[i].Name
					},
					fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
						if i == -1 {
							return "not found"
						}
						return "Name: " + tracks[i].Name + "\nArtist: " + tracks[i].Artist
					}),
					fuzzyfinder.WithMode(fuzzyfinder.ModeCaseSensitive),
				)
				if err != fuzzyfinder.ErrAbort {
					t.Fatalf("Find must return ErrAbort, but got '%s'", err)
				}

				return term.GetResult()
			})
		})
	}
}

func TestFind_enter(t *testing.T) {
	cases := map[string]struct {
		events   []termbox.Event
		expected int
	}{
		"initial":                      {events: keys(termbox.KeyTab), expected: 0},
		"mode smart to case-sensitive": {events: runes("CHI"), expected: 7},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			c := c
			events := c.events

			f, term := fuzzyfinder.NewWithMockedTerminal()
			events = append(events, key(termbox.KeyEnter))
			term.SetEvents(events...)

			idx, err := f.Find(
				tracks,
				func(i int) string {
					return tracks[i].Name
				},
			)
			if err != nil {
				t.Fatalf("Find must not return an error, but got '%s'", err)
			}
			if idx != c.expected {
				t.Errorf("expected index: %d, but got %d", c.expected, idx)
			}
		})
	}
}

func TestFind_error(t *testing.T) {
	t.Run("not a slice", func(t *testing.T) {
		f := fuzzyfinder.New()
		_, err := f.Find("", func(i int) string { return "" })
		if err == nil {
			t.Error("Find must return an error, but got nil")
		}
	})

	t.Run("itemFunc is nil", func(t *testing.T) {
		f := fuzzyfinder.New()
		_, err := f.Find([]string{}, nil)
		if err == nil {
			t.Error("Find must return an error, but got nil")
		}
	})
}

func TestFindMulti(t *testing.T) {
	cases := map[string]struct {
		events   []termbox.Event
		expected []int
		abort    bool
	}{
		"input glow":                          {events: runes("glow"), expected: []int{0}},
		"select two items":                    {events: keys(termbox.KeyTab, termbox.KeyArrowUp, termbox.KeyTab), expected: []int{0, 1}},
		"select two items with another order": {events: keys(termbox.KeyArrowUp, termbox.KeyTab, termbox.KeyTab), expected: []int{1, 0}},
		"toggle":                              {events: keys(termbox.KeyTab, termbox.KeyTab), expected: []int{0}},
		"empty result":                        {events: append(runes("ffffffffffffff")), abort: true},
		"resize window":                       {events: []termbox.Event{termbox.Event{Type: termbox.EventResize}}, expected: []int{0}},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			c := c
			events := c.events

			f, term := fuzzyfinder.NewWithMockedTerminal()
			events = append(events, key(termbox.KeyEnter))
			term.SetEvents(events...)

			idxs, err := f.FindMulti(
				tracks,
				func(i int) string {
					return tracks[i].Name
				},
				fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
					if i == -1 {
						return "not found"
					}
					return "Name: " + tracks[i].Name + "\nArtist: " + tracks[i].Artist
				}),
			)
			if c.abort {
				if err != fuzzyfinder.ErrAbort {
					t.Fatalf("Find must return ErrAbort, but got '%s'", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Find must not return an error, but got '%s'", err)
			}
			expectedSelectedNum := len(c.expected)
			if n := len(idxs); n != expectedSelectedNum {
				t.Errorf("expected the number of selected items is %d, but actual %d", expectedSelectedNum, n)
			}
		})
	}
}

func BenchmarkFind(b *testing.B) {
	f, term := fuzzyfinder.NewWithMockedTerminal()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		term.SetEvents(append(runes("adrele!!"), key(termbox.KeyEsc))...)
		f.Find(
			tracks,
			func(i int) string {
				return tracks[i].Name
			},
			fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
				if i == -1 {
					return "not found"
				}
				return "Name: " + tracks[i].Name + "\nArtist: " + tracks[i].Artist
			}),
		)
	}
}

func runes(s string) []termbox.Event {
	r := []rune(s)
	e := make([]termbox.Event, 0, len(r))
	for _, r := range r {
		e = append(e, ch(r))
	}
	return e
}

func ch(r rune) termbox.Event {
	if r == ' ' {
		return key(termbox.KeySpace)
	}
	return termbox.Event{
		Type: termbox.EventKey,
		Ch:   r,
	}
}

func key(key termbox.Key) termbox.Event {
	return termbox.Event{
		Type: termbox.EventKey,
		Key:  key,
	}
}

func keys(keys ...termbox.Key) []termbox.Event {
	k := make([]termbox.Event, 0, len(keys))
	for _, _key := range keys {
		k = append(k, key(_key))
	}
	return k
}
