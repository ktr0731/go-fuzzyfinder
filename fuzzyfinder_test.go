package fuzzyfinder_test

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ktr0731/iv/fuzzyfinder"
	"github.com/ktr0731/iv/logger"
	termbox "github.com/nsf/termbox-go"
)

var (
	update = flag.Bool("update", false, "update golden files")
	real   = flag.Bool("real", false, "display the actual layout to the terminal")
)

func init() {
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
		logger.Printf("golden updated: %s", fname)
		return
	}

	// Load the golden file.
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatalf("failed to load a golden file: %s", err)
	}
	expected := string(b)

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("wrong result: \n%s", diff)
	}
}

type track struct {
	Name   string
	Artist string
	Album  string
}

func newMockedTerminal() *fuzzyfinder.TerminalMock {
	mock := fuzzyfinder.UseMockedTerminal()
	w, h := 60, 10 // A normally value.
	mock.SetSize(w, h)
	return mock
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
		"arrow up-down":    {[]termbox.Event{key(termbox.KeyArrowUp)}},
		"arrow left-right": {append(runes("ゆるふわ樹海ガール"), key(termbox.KeyArrowLeft), key(termbox.KeyArrowLeft), key(termbox.KeyArrowLeft), key(termbox.KeyArrowRight))},
		"backspace":        {append(runes("adrenaline!!! -TV Ver.-"), key(termbox.KeyBackspace), key(termbox.KeyBackspace))},
		"backspace2":       {append(runes("オレンジ"), key(termbox.KeyBackspace2), key(termbox.KeyBackspace2))},
		"delete":           {append(runes("オレンジ"), key(termbox.KeyCtrlA), key(termbox.KeyDelete))},
		"ctrl-e":           {append(runes("恋をしたのは"), key(termbox.KeyCtrlA), key(termbox.KeyCtrlE))},
		"ctrl-w":           {append(runes("ハロ / ハワユ"), key(termbox.KeyCtrlW))},
		"ctrl-u":           {append(runes("恋をしたのは"), key(termbox.KeyArrowLeft), key(termbox.KeyCtrlU), key(termbox.KeyArrowRight))},
		"long item":        {[]termbox.Event{key(termbox.KeyArrowUp), key(termbox.KeyArrowUp), key(termbox.KeyArrowUp)}},
		"paging":           {[]termbox.Event{key(termbox.KeyArrowUp), key(termbox.KeyArrowUp), key(termbox.KeyArrowUp), key(termbox.KeyArrowUp), key(termbox.KeyArrowUp), key(termbox.KeyArrowUp), key(termbox.KeyArrowUp), key(termbox.KeyArrowUp)}},
		"multi":            {[]termbox.Event{key(termbox.KeyTab), key(termbox.KeyArrowUp), key(termbox.KeyTab)}},
		"backspace doesnt change x if cursorX is 0": {append(runes("a"), key(termbox.KeyCtrlA), key(termbox.KeyBackspace), key(termbox.KeyCtrlF))},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			c := c
			events := c.events

			term := newMockedTerminal()
			events = append(events, key(termbox.KeyEsc))
			term.SetEvents(events...)

			assertWithGolden(t, func(t *testing.T) string {
				_, err := fuzzyfinder.FindMulti(
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
				if err != fuzzyfinder.ErrAbort {
					t.Fatalf("Find must return ErrAbort, but got '%s'", err)
				}

				return term.GetResult()
			})
		})
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
