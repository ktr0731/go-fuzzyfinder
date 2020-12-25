package fuzzyfinder_test

import (
	"fmt"
	"io/ioutil"

	"github.com/gdamore/tcell/v2"
	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
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
	idx, _ := fuzzyfinder.Find(
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
			// the length of s, we return the name directly.
			if width < len([]rune(s)) {
				return slice[i].name
			}
			return s
		}))
	fmt.Println(slice[idx]) // The selected item.
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
	// Initialize a mocked terminal.
	term := fuzzyfinder.UseMockedTerminalV2()
	keys := "foo"
	for _, r := range keys {
		term.InjectKey(tcell.KeyRune, r, tcell.ModNone)
	}
	term.InjectKey(tcell.KeyEsc, rune(tcell.KeyEsc), tcell.ModNone)

	slice := []string{"foo", "bar", "baz"}
	_, _ = fuzzyfinder.Find(slice, func(i int) string { return slice[i] })

	// Write out the execution result to a temp file.
	// We can test it by the golden files testing pattern.
	//
	// See https://speakerdeck.com/mitchellh/advanced-testing-with-go?slide=19
	result := term.GetResult()
	_ = ioutil.WriteFile("ui.out", []byte(result), 0600)
}
