package fuzzyfinder

import (
	"github.com/gdamore/tcell/v2"
)

func New() *finder {
	return &finder{}
}

func NewWithMockedTerminal() (*finder, *TerminalMock) {
	f := New()
	m := f.UseMockedTerminal()
	w, h := 60, 10 // A normally value.
	m.Screen().(tcell.SimulationScreen).SetSize(w, h)
	return f, m
}
