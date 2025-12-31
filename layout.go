package fuzzyfinder

import "github.com/pkg/errors"

// rect represents a rectangular area with position and dimensions.
type rect struct {
	x      int
	y      int
	width  int
	height int
}

// splitVertical splits a rectangle vertically at the given ratio (0.0 to 1.0).
// Returns left and right rectangles.
func (r rect) splitVertical(ratio float64) (left, right rect) {
	if ratio <= 0 || ratio >= 1 {
		ratio = 0.5
	}

	splitX := int(float64(r.width) * ratio)

	left = rect{
		x:      r.x,
		y:      r.y,
		width:  splitX,
		height: r.height,
	}

	right = rect{
		x:      r.x + splitX,
		y:      r.y,
		width:  r.width - splitX,
		height: r.height,
	}

	return left, right
}

// inset creates a new rectangle by shrinking the current one by the given margin
// on all sides.
func (r rect) inset(margin int) rect {
	return rect{
		x:      r.x + margin,
		y:      r.y + margin,
		width:  r.width - 2*margin,
		height: r.height - 2*margin,
	}
}

// isEmpty returns true if the rectangle has zero or negative dimensions.
func (r rect) isEmpty() bool {
	return r.width <= 0 || r.height <= 0
}

// Layout represents the complete layout of the fuzzy finder interface.
// All coordinates are in terminal-absolute space (0,0 = top-left of screen).
type Layout struct {
	// terminal is the full terminal area
	terminal rect

	// content is the drawable area (may be smaller than terminal due to height/width options)
	content rect

	// border is the outer border area (same as content if border is enabled)
	border rect

	// innerContent is the area inside the border
	innerContent rect

	// list is the area for displaying matched items
	list rect

	// preview is the area for the preview panel (empty if no preview)
	preview rect

	// prompt is the area for the prompt line
	prompt rect

	// header is the area for the header line (empty if no header)
	header rect

	// numberLine is the area showing match count
	numberLine rect

	// items is the area for displaying the item list
	items rect

	// hasBorder indicates if a border is drawn
	hasBorder bool

	// hasPreview indicates if a preview panel is shown
	hasPreview bool

	// hasHeader indicates if a header is shown
	hasHeader bool
}

// validate checks if the layout has sufficient space for rendering.
func (l Layout) validate() error {
	if l.content.isEmpty() {
		return errors.New("content area is too small")
	}

	if l.innerContent.isEmpty() {
		return errors.New("inner content area is too small (border may be too large)")
	}

	if l.list.isEmpty() {
		return errors.New("list area is too small")
	}

	if l.items.isEmpty() {
		return errors.New("items area is too small (need at least 1 line for items)")
	}

	if l.list.width < 10 {
		return errors.New("terminal is too narrow (need at least 10 columns for list)")
	}

	if l.hasPreview && l.preview.isEmpty() {
		return errors.New("preview area is too small")
	}

	if l.hasPreview && l.preview.width < 10 {
		return errors.New("preview area is too narrow (need at least 10 columns)")
	}

	return nil
}

// computeLayout calculates the complete layout based on terminal size and options.
func (f *finder) computeLayout() (Layout, error) {
	layout := Layout{}

	// Get terminal size
	termWidth, termHeight := f.term.Size()
	layout.terminal = rect{x: 0, y: 0, width: termWidth, height: termHeight}

	// Apply width and height constraints
	contentWidth := termWidth
	contentHeight := termHeight

	if f.opt.width > 0 && f.opt.width < termWidth {
		contentWidth = f.opt.width
	}

	if f.opt.height > 0 && f.opt.height < termHeight {
		contentHeight = f.opt.height
	}

	// Position content area
	// Horizontally: center if width is constrained
	// Vertically: bottom if height is constrained (like fzf), otherwise use full terminal
	contentX := (termWidth - contentWidth) / 2
	contentY := 0
	if f.opt.height > 0 && f.opt.height < termHeight {
		// Position at bottom when height is explicitly set
		contentY = termHeight - contentHeight
	}

	layout.content = rect{x: contentX, y: contentY, width: contentWidth, height: contentHeight}
	layout.border = layout.content

	// Calculate inner content (accounting for border)
	layout.hasBorder = f.opt.border
	if layout.hasBorder {
		layout.innerContent = layout.content.inset(1)
	} else {
		layout.innerContent = layout.content
	}

	// Split horizontally for preview if needed
	layout.hasPreview = f.opt.previewFunc != nil
	if layout.hasPreview {
		layout.list, layout.preview = layout.innerContent.splitVertical(0.5)
	} else {
		layout.list = layout.innerContent
		layout.preview = rect{} // empty
	}

	// Now layout components vertically within the list area
	// Working from bottom to top (as the UI draws bottom-to-top)

	currentY := layout.list.y + layout.list.height - 1
	availableHeight := layout.list.height

	// Prompt line (always present, at bottom)
	if availableHeight < 1 {
		return layout, errors.New("insufficient height for prompt line")
	}
	layout.prompt = rect{
		x:      layout.list.x,
		y:      currentY,
		width:  layout.list.width,
		height: 1,
	}
	currentY--
	availableHeight--

	// Header line (optional)
	layout.hasHeader = len(f.opt.header) > 0
	if layout.hasHeader {
		if availableHeight < 1 {
			return layout, errors.New("insufficient height for header line")
		}
		layout.header = rect{
			x:      layout.list.x,
			y:      currentY,
			width:  layout.list.width,
			height: 1,
		}
		currentY--
		availableHeight--
	}

	// Number line (always present)
	if availableHeight < 1 {
		return layout, errors.New("insufficient height for number line")
	}
	layout.numberLine = rect{
		x:      layout.list.x,
		y:      currentY,
		width:  layout.list.width,
		height: 1,
	}
	currentY--
	availableHeight--

	// Items area (remaining space)
	if availableHeight < 1 {
		return layout, errors.New("insufficient height for items")
	}
	layout.items = rect{
		x:      layout.list.x,
		y:      layout.list.y,
		width:  layout.list.width,
		height: availableHeight,
	}

	return layout, nil
}
