package fuzzyfinder_test

import (
	"errors"
	"testing"

	"github.com/gdamore/tcell/v2"
	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
)

func FuzzPreviewWindow(f *testing.F) {
	slice := []string{"foo"}

	f.Add("Lorem ipsum dolor sit amet, consectetur adipiscing elit")
	f.Add("Sed eget dui libero.\nVivamus tempus, magna nec mollis convallis, ipsum justo tincidunt ligula, ut varius est mi id nisl.\nMorbi commodo turpis risus, nec vehicula leo auctor sit amet.\nUt imperdiet suscipit massa ac vehicula.\nInterdum et malesuada fames ac ante ipsum primis in faucibus.\nPraesent ligula orci, facilisis pulvinar varius eget, iaculis in erat.\nProin pellentesque arcu sed nisl consectetur tristique.\nQuisque tempus blandit dignissim.\nPhasellus dignissim sollicitudin mauris, sed gravida arcu luctus tincidunt.\nNunc rhoncus sed eros vel molestie.\nAenean sodales tortor eu libero rutrum, et lobortis orci scelerisque.\nPraesent sollicitudin, nunc ut consequat commodo, risus velit consectetur nibh, quis pretium nunc elit et erat.")
	f.Add("foo\x1b[31;1;44;0;90;105;38;5;12;48;5;226;38;2;10;20;30;48;2;200;100;50mbar")

	f.Fuzz(func(t *testing.T, s string) {
		finder, term := fuzzyfinder.NewWithMockedTerminal()
		events := []tcell.Event{key(input{tcell.KeyEsc, rune(tcell.KeyEsc), tcell.ModNone})}
		term.SetEventsV2(events...)

		_, err := finder.Find(
			slice,
			func(int) string { return slice[0] },
			fuzzyfinder.WithPreviewWindow(func(i, width, height int) string { return s }),
		)
		if !errors.Is(err, fuzzyfinder.ErrAbort) {
			t.Fatalf("Find must return ErrAbort, but got '%s'", err)
		}
	})
}
