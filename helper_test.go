package fuzzyfinder

func NewWithMockedTerminal() (*finder, *TerminalMock) {
	f := &finder{}
	m := f.UseMockedTerminal()
	w, h := 60, 10 // A normally value.
	m.SetSize(w, h)
	return f, m
}
