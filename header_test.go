package fuzzyfinder_test

import (
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
)

func TestHeaderOverlap(t *testing.T) {
	f, term := fuzzyfinder.NewWithMockedTerminal()

	// Set a small size
	term.SetSize(20, 5)

	items := []string{"item1", "item2", "item3"}
	header := "MY HEADER"

	// Event to close the finder immediately
	term.SetEventsV2(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))

	_, err := f.Find(
		items,
		func(i int) string { return items[i] },
		fuzzyfinder.WithHeader(header),
	)
	if err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	res := term.GetResult()

	// Check if Header is present
	if !strings.Contains(res, header) {
		t.Errorf("Result should contain header %q, but got:\n%s", header, res)
	}

	// Check if the first item is present
	// The first item "item1" should be visible.
	if !strings.Contains(res, "item1") {
		t.Errorf("Result should contain first item 'item1', but got:\n%s", res)
	}

	// More specific check: Header should be at the top (approximately)
	// and item should not overlap it.
	// Since GetResult returns a string with ANSI codes and newlines, we can split by newline.
	lines := strings.Split(res, "\n")
	
	// Helper to strip ANSI codes for easier checking (very basic stripping)
	strip := func(s string) string {
		var ret strings.Builder
		inEsc := false
		for _, r := range s {
			if r == '\x1b' {
				inEsc = true
				continue
			}
			if inEsc {
				if r == 'm' {
					inEsc = false
				}
				continue
			}
			ret.WriteRune(r)
		}
		return ret.String()
	}

	// We expect the header to be on the first line (index 0)
	// Note: terminal mock dump might have empty lines or reset codes at the start.
	// Let's find the line with the header.
	headerLineIdx := -1
	itemLineIdx := -1

	for i, line := range lines {
		clean := strip(line)
		if strings.Contains(clean, header) {
			headerLineIdx = i
		}
		// matched item has "> " prefix usually, or just check content
		if strings.Contains(clean, "item1") {
			itemLineIdx = i
		}
	}

	if headerLineIdx == -1 {
		t.Fatal("Could not find header line")
	}
	if itemLineIdx == -1 {
		t.Fatal("Could not find item line")
	}

	if headerLineIdx == itemLineIdx {
		t.Errorf("Header and First Item are on the same line (%d). Overlap detected!", headerLineIdx)
	}

	if headerLineIdx > itemLineIdx {
		t.Errorf("Header is below the item? Header line: %d, Item line: %d", headerLineIdx, itemLineIdx)
	}
}
