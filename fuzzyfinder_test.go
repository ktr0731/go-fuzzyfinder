package fuzzyfinder_test

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
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
		events []tcell.Event
	}{
		"initial":    {},
		"input lo":   {runes("lo")},
		"input glow": {runes("glow")},
		"arrow up-down": {keys([]input{
			{tcell.KeyUp, rune(tcell.KeyUp), tcell.ModNone},
			{tcell.KeyUp, rune(tcell.KeyUp), tcell.ModNone},
			{tcell.KeyDown, rune(tcell.KeyDown), tcell.ModNone},
		}...)},
		"arrow left-right": {append(runes("ゆるふわ樹海"), keys([]input{
			{tcell.KeyLeft, rune(tcell.KeyLeft), tcell.ModNone},
			{tcell.KeyLeft, rune(tcell.KeyLeft), tcell.ModNone},
			{tcell.KeyRight, rune(tcell.KeyRight), tcell.ModNone},
		}...)...)},
		"backspace": {append(runes("adr .-"), keys([]input{
			{tcell.KeyBackspace, rune(tcell.KeyBackspace), tcell.ModNone},
			{tcell.KeyBackspace, rune(tcell.KeyBackspace), tcell.ModNone},
		}...)...)},
		"backspace empty": {keys(input{tcell.KeyBackspace2, rune(tcell.KeyBackspace2), tcell.ModNone})},
		"backspace2": {append(runes("オレンジ"), keys([]input{
			{tcell.KeyBackspace2, rune(tcell.KeyBackspace2), tcell.ModNone},
			{tcell.KeyBackspace2, rune(tcell.KeyBackspace2), tcell.ModNone},
		}...)...)},
		"delete": {append(runes("オレンジ"), keys([]input{
			{tcell.KeyCtrlA, 'A', tcell.ModCtrl},
			{tcell.KeyDelete, rune(tcell.KeyDelete), tcell.ModNone},
		}...)...)},
		"delete empty": {keys([]input{
			{tcell.KeyCtrlA, 'A', tcell.ModCtrl},
			{tcell.KeyDelete, rune(tcell.KeyDelete), tcell.ModNone},
		}...)},
		"ctrl-e": {append(runes("恋をしたのは"), keys([]input{
			{tcell.KeyCtrlA, 'A', tcell.ModCtrl},
			{tcell.KeyCtrlE, 'E', tcell.ModCtrl},
		}...)...)},
		"ctrl-w":       {append(runes("ハロ / ハワユ"), keys(input{tcell.KeyCtrlW, 'W', tcell.ModCtrl})...)},
		"ctrl-w empty": {keys(input{tcell.KeyCtrlW, 'W', tcell.ModCtrl})},
		"ctrl-u": {append(runes("恋をしたのは"), keys([]input{
			{tcell.KeyLeft, rune(tcell.KeyLeft), tcell.ModNone},
			{tcell.KeyCtrlU, 'U', tcell.ModCtrl},
			{tcell.KeyRight, rune(tcell.KeyRight), tcell.ModNone},
		}...)...)},
		"long item": {keys([]input{
			{tcell.KeyUp, rune(tcell.KeyUp), tcell.ModNone},
			{tcell.KeyUp, rune(tcell.KeyUp), tcell.ModNone},
			{tcell.KeyUp, rune(tcell.KeyUp), tcell.ModNone},
		}...)},
		"paging": {keys([]input{
			{tcell.KeyUp, rune(tcell.KeyUp), tcell.ModNone},
			{tcell.KeyUp, rune(tcell.KeyUp), tcell.ModNone},
			{tcell.KeyUp, rune(tcell.KeyUp), tcell.ModNone},
			{tcell.KeyUp, rune(tcell.KeyUp), tcell.ModNone},
			{tcell.KeyUp, rune(tcell.KeyUp), tcell.ModNone},
			{tcell.KeyUp, rune(tcell.KeyUp), tcell.ModNone},
			{tcell.KeyUp, rune(tcell.KeyUp), tcell.ModNone},
			{tcell.KeyUp, rune(tcell.KeyUp), tcell.ModNone},
		}...)},
		"tab doesn't work": {keys(input{tcell.KeyTab, rune(tcell.KeyTab), tcell.ModNone})},
		"backspace doesnt change x if cursorX is 0": {append(runes("a"), keys([]input{
			{tcell.KeyCtrlA, 'A', tcell.ModCtrl},
			{tcell.KeyBackspace, rune(tcell.KeyBackspace), tcell.ModNone},
			{tcell.KeyCtrlF, 'F', tcell.ModCtrl},
		}...)...)},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			events := c.events

			f, term := fuzzyfinder.NewWithMockedTerminal()
			events = append(events, key(input{tcell.KeyEsc, rune(tcell.KeyEsc), tcell.ModNone}))
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

				res := term.GetResult()
				term.Screen().Fini()
				return res
			})
		})
	}
}

func TestFind_hotReload(t *testing.T) {
	f, term := fuzzyfinder.NewWithMockedTerminal()
	events := append(runes("adrena"), keys(input{tcell.KeyEsc, rune(tcell.KeyEsc), tcell.ModNone})...)
	term.SetEvents(events...)

	var mu sync.Mutex
	assertWithGolden(t, func(t *testing.T) string {
		_, err := f.Find(
			&tracks,
			func(i int) string {
				mu.Lock()
				defer mu.Unlock()
				return tracks[i].Name
			},
			fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
				// Hack, wait until updateItems is called.
				time.Sleep(50 * time.Millisecond)
				mu.Lock()
				defer mu.Unlock()
				if i == -1 {
					return "not found"
				}
				return "Name: " + tracks[i].Name + "\nArtist: " + tracks[i].Artist
			}),
			fuzzyfinder.WithMode(fuzzyfinder.ModeCaseSensitive),
			fuzzyfinder.WithHotReload(),
		)
		if err != fuzzyfinder.ErrAbort {
			t.Fatalf("Find must return ErrAbort, but got '%s'", err)
		}

		res := term.GetResult()
		term.Screen().Fini()
		return res
	})
}

func TestFind_enter(t *testing.T) {
	cases := map[string]struct {
		events   []tcell.Event
		expected int
	}{
		"initial": {events: keys(input{tcell.KeyTab, rune(tcell.KeyTab), tcell.ModNone}), expected: 0},
		"mode smart to case-sensitive": {events: append(runes("JI"), keys([]input{
			{tcell.KeyTab, rune(tcell.KeyTab), tcell.ModNone}, // tab earn time for filter
			{tcell.KeyTab, rune(tcell.KeyTab), tcell.ModNone},
			{tcell.KeyTab, rune(tcell.KeyTab), tcell.ModNone},
			{tcell.KeyTab, rune(tcell.KeyTab), tcell.ModNone},
			{tcell.KeyTab, rune(tcell.KeyTab), tcell.ModNone},
			{tcell.KeyTab, rune(tcell.KeyTab), tcell.ModNone},
			{tcell.KeyTab, rune(tcell.KeyTab), tcell.ModNone},
		}...)...), expected: 7},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			c := c
			events := c.events

			f, term := fuzzyfinder.NewWithMockedTerminal()
			events = append(events, key(input{tcell.KeyEnter, rune(tcell.KeyEnter), tcell.ModNone}))
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
		events   []tcell.Event
		expected []int
		abort    bool
	}{
		"input glow": {events: runes("glow"), expected: []int{0}},
		"select two items": {events: keys([]input{
			{tcell.KeyTab, rune(tcell.KeyTab), tcell.ModNone},
			{tcell.KeyUp, rune(tcell.KeyUp), tcell.ModNone},
			{tcell.KeyTab, rune(tcell.KeyTab), tcell.ModNone},
		}...), expected: []int{0, 1}},
		"select two items with another order": {events: keys([]input{
			{tcell.KeyUp, rune(tcell.KeyUp), tcell.ModNone},
			{tcell.KeyTab, rune(tcell.KeyTab), tcell.ModNone},
			{tcell.KeyTab, rune(tcell.KeyTab), tcell.ModNone},
		}...), expected: []int{1, 0}},
		"toggle": {events: keys([]input{
			{tcell.KeyTab, rune(tcell.KeyTab), tcell.ModNone},
			{tcell.KeyTab, rune(tcell.KeyTab), tcell.ModNone},
		}...), expected: []int{0}},
		"empty result": {events: append(runes("fffffff")), abort: true},
		"resize window": {events: []tcell.Event{
			tcell.NewEventResize(10, 10),
		}, expected: []int{0}},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			c := c
			events := c.events

			f, term := fuzzyfinder.NewWithMockedTerminal()
			events = append(events, key(input{tcell.KeyEnter, rune(tcell.KeyEnter), tcell.ModNone}))
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
	b.Run("normal", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			f, term := fuzzyfinder.NewWithMockedTerminal()
			term.SetEvents(append(runes("adrele!!"), key(input{tcell.KeyEsc, rune(tcell.KeyEsc), tcell.ModNone}))...)
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
	})

	b.Run("hotreload", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			f, term := fuzzyfinder.NewWithMockedTerminal()
			term.SetEvents(append(runes("adrele!!"), key(input{tcell.KeyEsc, rune(tcell.KeyEsc), tcell.ModNone}))...)
			f.Find(
				&tracks,
				func(i int) string {
					return tracks[i].Name
				},
				fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
					if i == -1 {
						return "not found"
					}
					return "Name: " + tracks[i].Name + "\nArtist: " + tracks[i].Artist
				}),
				fuzzyfinder.WithHotReload(),
			)
		}
	})
}

func runes(s string) []tcell.Event {
	r := []rune(s)
	e := make([]tcell.Event, 0, len(r))
	for _, r := range r {
		e = append(e, ch(r))
	}
	return e
}

func ch(r rune) tcell.Event {
	return key(input{tcell.KeyRune, r, tcell.ModNone})
}

func key(input input) tcell.Event {
	return tcell.NewEventKey(input.key, input.ch, input.mod)
}

func keys(inputs ...input) []tcell.Event {
	k := make([]tcell.Event, 0, len(inputs))
	for _, in := range inputs {
		k = append(k, key(in))
	}
	return k
}

type input struct {
	key tcell.Key
	ch  rune
	mod tcell.ModMask
}
