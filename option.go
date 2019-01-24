package fuzzyfinder

type opt struct {
	previewFunc func(i, width, height int) string
	multi       bool
}

type option func(*opt)

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
