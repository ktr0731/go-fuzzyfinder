package fuzzyfinder

import "github.com/gdamore/tcell/v2"

func New() *finder {
	return &finder{}
}

func NewWithMockedTerminal() (*finder, *TerminalMock) {
	eventsChan := make(chan tcell.Event, 10)
	quitChan := make(chan struct{})

	f := New()
	f.termEventsChan = eventsChan
	f.termQuitChan = quitChan

	m := f.UseMockedTerminalV2()
	go m.ChannelEvents(eventsChan, quitChan)

	w, h := 60, 10 // A normally value.
	m.SetSize(w, h)
	return f, m
}
