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
