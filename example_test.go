package fuzzyfinder_test

import (
	"fmt"
	"io/ioutil"

	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
	"github.com/nsf/termbox-go"
)

func ExampleFind() {
	slice := []struct {
		id   string
		name string
	}{
		{"id1", "foo"},
		{"id2", "bar"},
		{"id3", "baz"},
	}
	idx, _ := fuzzyfinder.Find(slice, func(i int) string {
		return fmt.Sprintf("[%s] %s", slice[i].id, slice[i].name)
	})
	fmt.Println(slice[idx]) // The selected item.
}

func ExampleFind_previewWindow() {
	slice := []struct {
		id   string
		name string
	}{
		{"id1", "foo"},
		{"id2", "bar"},
		{"id3", "baz"},
	}
	fuzzyfinder.Find(
		slice,
		func(i int) string {
			return fmt.Sprintf("[%s] %s", slice[i].id, slice[i].name)
		},
		fuzzyfinder.WithPreviewWindow(func(i, width, _ int) string {
			if i == -1 {
				return "no results"
			}
			s := fmt.Sprintf("%s is selected", slice[i].name)
			// As an example of using width, if the window width is less than
			// the length of s, we returns the name directly.
			if width < len([]rune(s)) {
				return slice[i].name
			}
			return s
		}))
}

func ExampleFindMulti() {
	slice := []struct {
		id   string
		name string
	}{
		{"id1", "foo"},
		{"id2", "bar"},
		{"id3", "baz"},
	}
	idxs, _ := fuzzyfinder.FindMulti(slice, func(i int) string {
		return fmt.Sprintf("[%s] %s", slice[i].id, slice[i].name)
	})
	for _, idx := range idxs {
		fmt.Println(slice[idx])
	}
}

func ExampleTerminalMock() {
	keys := func(str string) []termbox.Event {
		s := []rune(str)
		e := make([]termbox.Event, 0, len(s))
		for _, r := range s {
			e = append(e, termbox.Event{Type: termbox.EventKey, Ch: r})
		}
		return e
	}

	// Initialize a mocked terminal.
	term := fuzzyfinder.UseMockedTerminal()
	// Set the window size and events.
	term.SetSize(60, 10)
	term.SetEvents(append(
		keys("foo"),
		termbox.Event{Type: termbox.EventKey, Key: termbox.KeyEsc})...)

	// Call fuzzyfinder.Find.
	slice := []string{"foo", "bar", "baz"}
	fuzzyfinder.Find(slice, func(i int) string { return slice[i] })

	// Write out the execution result to a temp file.
	// We can test it by the golden files testing pattern.
	//
	// See https://speakerdeck.com/mitchellh/advanced-testing-with-go?slide=19
	res := term.GetResult()
	ioutil.WriteFile("ui.out", []byte(res), 0644)
}
