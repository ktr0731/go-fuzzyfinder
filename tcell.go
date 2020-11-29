package fuzzyfinder

import (
	"github.com/gdamore/tcell/v2"
)

type terminal interface {
	Screen() tcell.Screen
}

// termImpl is the implementation for termbox-go.
type termImpl struct {
	screen tcell.Screen
}

func (t *termImpl) Screen() tcell.Screen {
	return t.screen
}
