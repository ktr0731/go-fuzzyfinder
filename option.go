package fuzzyfinder

type opt struct {
	mode        mode
	previewFunc func(i, width, height int) string
	multi       bool
}

type mode int

const (
	// ModeSmart enables a smart matching. It is the default matching mode.
	// At the beginning, matching mode is ModeCaseInsensitive, but it switches
	// over to ModeCaseSensitive if an upper case character is inputted.
	ModeSmart mode = iota
	// ModeCaseSensitive enables a case-sensitive matching.
	ModeCaseSensitive
	// ModeCaseInsensitive enables a case-insensitive matching.
	ModeCaseInsensitive
)

type option func(*opt)

// WithMode specifies a matching mode.
func WithMode(m mode) option {
	return func(o *opt) {
		o.mode = m
	}
}

// WithPreviewWindow enables to display a preview for the selected item.
// the argument f receives i, width and height. i is the same as Find's one.
// width and height are the size of the terminal. You can use these to adjust
// a preview content. If there is no selected item, previewFunc passes -1 to
// previewFunc.
func WithPreviewWindow(f func(i, width, height int) string) option {
	return func(o *opt) {
		o.previewFunc = f
	}
}

// withMulti enables to select multiple items by tab key.
func withMulti() option {
	return func(o *opt) {
		o.multi = true
	}
}
