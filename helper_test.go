package fuzzyfinder

func New() *finder {
	return &finder{}
}

func NewWithMockedTerminal() (*finder, *TerminalMock) {
	f := New()
	m := f.UseMockedTerminal()
	w, h := 60, 10 // A normally value.
	m.SetSize(w, h)
	return f, m
}

func NewWithMockedTerminalV2() (*finder, *TerminalMock) {
	f := New()
	m := f.UseMockedTerminalV2()
	w, h := 60, 10 // A normally value.
	m.SetSizeV2(w, h)
	return f, m
}
