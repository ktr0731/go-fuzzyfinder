package strmatch

// Matched represents a result of FindAll.
type Matched struct {
	// Idx is the index of an item of the original slice which was used to
	// search matched strings.
	Idx int
	// Pos is the range of matched position.
	// [2]int represents a closed interval of a position.
	Pos [2]int
}

type opt optFunc
type optFunc func(*option)

// option represents available options and its default values.
type option struct {
	caseSensitive bool
}

// WithCaseSensitive enables a case sensitive searching.
func WithCaseSensitive() opt {
	return func(o *option) {
		o.caseSensitive = true
	}
}

// FindAll tries to find out sub-strings from slice that match the passed argument in.
func FindAll(in string, slice []string, opts ...opt) []Matched {
	var opt option
	for _, o := range opts {
		o(&opt)
	}
	return match(in, slice, opt)
}
