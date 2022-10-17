//go:build fuzz
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

	"github.com/gdamore/tcell/v2"

	fuzz "github.com/google/gofuzz"
	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
)

type fuzzKey struct {
	key  tcell.Key
	name string
}

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789一花二乃三玖四葉五月")
	tbkeys  = []tcell.Key{
		tcell.KeyCtrlA,
		tcell.KeyCtrlB,
		tcell.KeyCtrlE,
		tcell.KeyCtrlF,
		tcell.KeyBackspace,
		tcell.KeyTab,
		tcell.KeyCtrlJ,
		tcell.KeyCtrlK,
		tcell.KeyCtrlN,
		tcell.KeyCtrlP,
		tcell.KeyCtrlU,
		tcell.KeyCtrlW,
		tcell.KeyBackspace2,
		tcell.KeyUp,
		tcell.KeyDown,
		tcell.KeyLeft,
		tcell.KeyRight,
	}
	keyMap = map[tcell.Key]string{
		tcell.KeyCtrlA:      "A",
		tcell.KeyCtrlB:      "B",
		tcell.KeyCtrlE:      "E",
		tcell.KeyCtrlF:      "F",
		tcell.KeyBackspace:  "backspace",
		tcell.KeyTab:        "tab",
		tcell.KeyCtrlJ:      "J",
		tcell.KeyCtrlK:      "K",
		tcell.KeyCtrlN:      "N",
		tcell.KeyCtrlP:      "P",
		tcell.KeyCtrlU:      "U",
		tcell.KeyCtrlW:      "W",
		tcell.KeyBackspace2: "backspace2",
		tcell.KeyUp:         "up",
		tcell.KeyDown:       "down",
		tcell.KeyLeft:       "left",
		tcell.KeyRight:      "right",
	}
)

var (
	out       = flag.String("fuzzout", "fuzz.out", "fuzzing error cases")
	hotReload = flag.Bool("hotreload", false, "enable hot-reloading")
	numCases  = flag.Int("numCases", 30, "number of test cases")
	numEvents = flag.Int("numEvents", 10, "number of events")
)

// TestFuzz executes fuzzing tests.
//
// Example:
//
//   go test -tags fuzz -run TestFuzz -numCases 10 -numEvents 10
//
func TestFuzz(t *testing.T) {
	f, err := os.Create(*out)
	if err != nil {
		t.Fatalf("failed to create a fuzzing output file: %s", err)
	}
	defer f.Close()

	fuzz := fuzz.New()

	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	for i := 0; i < rand.Intn(*numCases)+10; i++ {
		// number of events in tcell.SimulationScreen is limited 10
		n := rand.Intn(min(*numEvents, 10))
		events := make([]tcell.Event, n)
		for i := 0; i < n; i++ {
			if rand.Intn(10) > 3 {
				events[i] = ch(letters[rand.Intn(len(letters)-1)])
			} else {
				k := tbkeys[rand.Intn(len(tbkeys)-1)]
				events[i] = key(input{k, rune(k), tcell.ModNone})
			}
		}

		var name string
		for _, e := range events {
			if e.(*tcell.EventKey).Rune() != 0 {
				name += string(e.(*tcell.EventKey).Rune())
			} else {
				name += "[" + keyMap[e.(*tcell.EventKey).Key()] + "]"
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
			events = append(events, key(input{tcell.KeyEsc, rune(tcell.KeyEsc), tcell.ModNone}))

			term.SetEventsV2(events...)

			var (
				iface     interface{}
				promptStr string
				header    string
			)
			fuzz.Fuzz(&promptStr)
			fuzz.Fuzz(&header)
			opts := []fuzzyfinder.Option{
				fuzzyfinder.WithPromptString(promptStr),
				fuzzyfinder.WithHeader(header),
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
