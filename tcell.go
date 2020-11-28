package fuzzyfinder

import (
	"github.com/gdamore/tcell/v2"
)

type terminal interface {
	screen() tcell.Screen
}

// termImpl is the implementation for termbox-go.
type termImpl struct {
	s tcell.Screen
}

func (t *termImpl) screen() tcell.Screen {
	return t.s
}
