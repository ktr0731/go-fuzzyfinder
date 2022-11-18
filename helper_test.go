package fuzzyfinder

import "github.com/gdamore/tcell/v2"

func New() *finder {
	return &finder{}
}

func NewWithMockedTerminal() (*finder, *TerminalMock) {
	eventsChan := make(chan tcell.Event, 10)

	f := New()
	f.termEventsChan = eventsChan

	m := f.UseMockedTerminalV2()
	go m.ChannelEvents(eventsChan, nil)

	w, h := 60, 10 // A normally value.
	m.SetSize(w, h)
	return f, m
}
