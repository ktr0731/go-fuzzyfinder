// +build fuzz

package fuzzyfinder_test

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"

	fuzz "github.com/google/gofuzz"
	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
	"github.com/nsf/termbox-go"
)

type fuzzKey struct {
	key  termbox.Key
	name string
}

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789一花二乃三玖四葉五月")
	tbkeys  = []termbox.Key{
		termbox.KeyCtrlA,
		termbox.KeyCtrlB,
		termbox.KeyCtrlE,
		termbox.KeyCtrlF,
		termbox.KeyBackspace,
		termbox.KeyTab,
		termbox.KeyCtrlJ,
		termbox.KeyCtrlK,
		termbox.KeyCtrlN,
		termbox.KeyCtrlP,
		termbox.KeyCtrlU,
		termbox.KeyCtrlW,
		termbox.KeySpace,
		termbox.KeyBackspace2,
		termbox.KeyArrowUp,
		termbox.KeyArrowDown,
		termbox.KeyArrowLeft,
		termbox.KeyArrowRight,
	}
	keyMap = map[termbox.Key]string{
		termbox.KeyCtrlA:      "A",
		termbox.KeyCtrlB:      "B",
		termbox.KeyCtrlE:      "E",
		termbox.KeyCtrlF:      "F",
		termbox.KeyBackspace:  "backspace",
		termbox.KeyTab:        "tab",
		termbox.KeyCtrlJ:      "J",
		termbox.KeyCtrlK:      "K",
		termbox.KeyCtrlN:      "N",
		termbox.KeyCtrlP:      "P",
		termbox.KeyCtrlU:      "U",
		termbox.KeyCtrlW:      "W",
		termbox.KeySpace:      "space",
		termbox.KeyBackspace2: "backspace2",
		termbox.KeyArrowUp:    "up",
		termbox.KeyArrowDown:  "down",
		termbox.KeyArrowLeft:  "left",
		termbox.KeyArrowRight: "right",
	}
)

var (
	out       = flag.String("fuzzout", "fuzz.out", "fuzzing error cases")
	hotReload = flag.Bool("hotreload", false, "enable hot-reloading")
	numCases  = flag.Int("numCases", 30, "number of test cases")
	numEvents = flag.Int("numEvents", 100, "number of events")
)

// TestFuzz executes fuzzing tests.
//
// Example:
//
//   go test -tags fuzz -run TestFuzz -numCases 1000 -numEvents 100
//
func TestFuzz(t *testing.T) {
	f, err := os.Create(*out)
	if err != nil {
		t.Fatalf("failed to create a fuzzing output file: %s", err)
	}
	defer f.Close()

	fuzz := fuzz.New()

	for i := 0; i < rand.Intn(*numCases)+10; i++ {
		n := rand.Intn(*numEvents) + 10
		events := make([]termbox.Event, n)
		for i := 0; i < n; i++ {
			if rand.Intn(10) > 3 {
				events[i] = ch(letters[rand.Intn(len(letters)-1)])
			} else {
				events[i] = key(tbkeys[rand.Intn(len(tbkeys)-1)])
			}
		}

		var name string
		for _, e := range events {
			if e.Key == termbox.KeySpace {
				name += " "
			} else if e.Ch != 0 {
				name += string(e.Ch)
			} else {
				name += "[" + keyMap[e.Key] + "]"
			}
		}

		t.Run(name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Fprintln(f, name)
					t.Errorf("panicked: %s", name)
				}
				return
			}()

			var mu sync.Mutex
			tracks := tracks

			f, term := fuzzyfinder.NewWithMockedTerminal()
			events = append(events, key(termbox.KeyEsc))
			term.SetEvents(events...)

			var (
				iface     interface{}
				promptStr string
			)
			fuzz.Fuzz(&promptStr)
			opts := []fuzzyfinder.Option{
				fuzzyfinder.WithPromptString(promptStr),
			}
			if *hotReload {
				iface = &tracks
				opts = append(opts, fuzzyfinder.WithHotReload())
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				go func() {
					for {
						select {
						case <-ctx.Done():
							return
						default:
							var t track
							fuzz.Fuzz(&t.Name)
							fuzz.Fuzz(&t.Artist)
							fuzz.Fuzz(&t.Album)
							mu.Lock()
							tracks = append(tracks, &t)
							mu.Unlock()
						}
					}
				}()
			} else {
				iface = tracks
			}

			_, err := f.Find(
				iface,
				func(i int) string {
					mu.Lock()
					defer mu.Unlock()
					return tracks[i].Name
				},
				append(
					opts,
					fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
						if i == -1 {
							return "not found"
						}
						mu.Lock()
						defer mu.Unlock()
						return "Name: " + tracks[i].Name + "\nArtist: " + tracks[i].Artist
					}),
				)...,
			)
			if err != fuzzyfinder.ErrAbort {
				t.Fatalf("Find must return ErrAbort, but got '%s'", err)
			}

		})
	}
}
