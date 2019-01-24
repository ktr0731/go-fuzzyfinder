package strmatch

// Matched represents a result of FindAll.
type Matched struct {
	// Idx is the index of an item of the original slice which was used to
	// search matched strings.
	Idx int
	// Pos is the collection of a matched position.
	// [2]int represents a closed interval of a position.
	Pos [][2]int
}

// algo represents an algorithm.
type algo func(in string, slice []string, opt option) []Matched

var (
	// AlgoNoFuzzy provides a non-fuzzy algorithm.
	// It matches only a continuous sub-string. For that reason, each
	// Matched which is returned by FindAll with AlgoNoFuzzy returns only
	// one position.
	AlgoNoFuzzy algo = noFuzzy
	// AlgoRegExp provides a fuzzy searching using regexp.
	// Each matched which is returned by FindAll with AlgoRegExp returns only
	// one position. It includes from the first rune of input to
	// the last rune of input.
	//
	// For example, we assume that the target string is "Twinkle Snow" and
	// input is "ink now". This algorithm returns "inkle Snow" as a matched
	// string.
	AlgoRegExp algo = regexpMatch
)

type opt optFunc
type optFunc func(*option)

// option represents available options and its default values.
type option struct {
	algo          algo
	caseSensitive bool
}

// WithAlgo changes the algorithm which is used to FindAll.
func WithAlgo(algo algo) opt {
	return func(o *option) {
		o.algo = algo
	}
}

// WithCaseSensitive enables a case sensitive searching.
func WithCaseSensitive() opt {
	return func(o *option) {
		o.caseSensitive = true
	}
}

// FindAll tries to find out sub-strings from slice that match the passed argument in.
// The default algorithm which is used to fuzzy searching is AlgoNoFuzzy.
func FindAll(in string, slice []string, opts ...opt) []Matched {
	var opt option
	for _, o := range opts {
		o(&opt)
	}
	return opt.algo(in, slice, opt)
}
