package fuzzyfinder

import "time"

func New() *finder {
	return &finder{}
}

func NewWithMockedTerminal() (*finder, *TerminalMock) {
	f := New()
	m := f.UseMockedTerminal()
	w, h := 60, 10 // A normally value.
	m.SetSize(w, h)
	m.sleepDuration = 500 * time.Microsecond
	return f, m
}
