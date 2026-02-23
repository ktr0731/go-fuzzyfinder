package fuzzyfinder_test

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
	"github.com/pkg/errors"
)

func TestFind_WithSearchItemFunc(t *testing.T) {
	type item struct {
		Name string
		ID   string
	}
	items := []item{
		{"Foo", "id-123"},
		{"Bar", "id-456"},
		{"Baz", "id-789"},
	}

	cases := map[string]struct {
		events []tcell.Event
	}{
		"search by hidden id": {
			events: append(
				append([]tcell.Event{key(input{tcell.KeyCtrlO, rune(tcell.KeyCtrlO), tcell.ModCtrl})}, runes("id-456")...),
				key(input{tcell.KeyEsc, rune(tcell.KeyEsc), tcell.ModNone})),
		},
		"search by hidden id partial": {
			events: append(
				append([]tcell.Event{key(input{tcell.KeyCtrlO, rune(tcell.KeyCtrlO), tcell.ModCtrl})}, runes("id-7")...),
				key(input{tcell.KeyEsc, rune(tcell.KeyEsc), tcell.ModNone})),
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			f, term := fuzzyfinder.NewWithMockedTerminal()
			term.SetEventsV2(c.events...)

			assertWithGolden(t, func(t *testing.T) string {
				_, err := f.Find(
					items,
					func(i int) string {
						return items[i].Name
					},
					fuzzyfinder.WithSearchItemFunc(func(i int) string {
						return items[i].ID
					}),
				)
				if !errors.Is(err, fuzzyfinder.ErrAbort) {
					t.Fatalf("Find must return ErrAbort, but got '%s'", err)
				}
				return term.GetResult()
			})
		})
	}
}
