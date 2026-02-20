package fuzzyfinder_test

import (
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
)

func TestScrollHorizontal(t *testing.T) {
	f, term := fuzzyfinder.NewWithMockedTerminal()

	// Use a small width to force scrolling
	// Width 10. Item starts at x=2. Visible width for item = 8.
	term.SetSize(10, 10)

	longItem := "0123456789ABCDEFGHIJ"
	items := []string{longItem}

	// Inject keys:
	// 1. Shift+Right -> Scroll right
	// 2. Shift+Right -> Scroll right again
	// 3. Enter -> Select
	//
	// To test Shift+Left, we can create another test case or add more steps.
	// Let's add Shift+Left to restore the view.
	events := []tcell.Event{
		tcell.NewEventKey(tcell.KeyRight, ' ', tcell.ModShift),
		tcell.NewEventKey(tcell.KeyRight, ' ', tcell.ModShift),
		// Check intermediate state? We can't easily with this setup as Find blocks until return.
		// So we just test right scroll first.
		tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone),
	}

	term.SetEventsV2(events...)

	_, err := f.Find(items, func(i int) string { return items[i] })
	if err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	res := term.GetResult()

	if !strings.Contains(res, "2345") {
		t.Errorf("Expected result to contain '2345' (scrolled content), but got:\n%s", res)
	}
	if strings.Contains(res, "01") {
		t.Errorf("Expected result NOT to contain '01' (scrolled out), but got:\n%s", res)
	}

	// Test Shift+Left
	term.SetSize(10, 10) // Reset size (though redundant)
	eventsLeft := []tcell.Event{
		tcell.NewEventKey(tcell.KeyRight, ' ', tcell.ModShift),
		tcell.NewEventKey(tcell.KeyRight, ' ', tcell.ModShift),
		tcell.NewEventKey(tcell.KeyLeft, ' ', tcell.ModShift),
		tcell.NewEventKey(tcell.KeyLeft, ' ', tcell.ModShift),
		tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone),
	}
	term.SetEventsV2(eventsLeft...)

	_, err = f.Find(items, func(i int) string { return items[i] })
	if err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	resLeft := term.GetResult()
	if !strings.Contains(resLeft, "01") {
		t.Errorf("Expected result to contain '01' (scrolled back), but got:\n%s", resLeft)
	}
}
