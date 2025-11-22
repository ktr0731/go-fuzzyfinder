package fuzzyfinder

import (
	"time"
	"unicode"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	runewidth "github.com/mattn/go-runewidth"
)

const (
	// maxUndoHistory is the maximum number of undo states to keep in history
	maxUndoHistory = 100
	// doubleClickThreshold is the maximum time between clicks to consider it a double-click
	doubleClickThreshold = 500 * time.Millisecond
)

// undoState represents a snapshot of the input state for undo/redo
type undoState struct {
	input          []rune
	x              int
	cursorX        int
	selectionStart int
	selectionEnd   int
}

// saveUndoState saves the current input state to the undo history
func (f *finder) saveUndoState() {
	// Truncate history if we're not at the end
	if f.state.undoIndex < len(f.state.undoHistory)-1 {
		f.state.undoHistory = f.state.undoHistory[:f.state.undoIndex+1]
	}

	// Save current state
	snapshot := undoState{
		input:          append([]rune(nil), f.state.input...),
		x:              f.state.x,
		cursorX:        f.state.cursorX,
		selectionStart: f.state.selectionStart,
		selectionEnd:   f.state.selectionEnd,
	}

	f.state.undoHistory = append(f.state.undoHistory, snapshot)
	f.state.undoIndex = len(f.state.undoHistory) - 1

	// Limit history size to maxUndoHistory entries
	if len(f.state.undoHistory) > maxUndoHistory {
		f.state.undoHistory = f.state.undoHistory[1:]
		f.state.undoIndex--
	}
}

// undo reverts to the previous state in the undo history
func (f *finder) undo() {
	// Need at least 2 states (current + one to revert to)
	if f.state.undoIndex < 1 || len(f.state.undoHistory) < 2 {
		return // Nothing to undo
	}

	f.state.undoIndex--
	snapshot := f.state.undoHistory[f.state.undoIndex]

	f.state.input = append([]rune(nil), snapshot.input...)
	f.state.x = snapshot.x
	f.state.cursorX = snapshot.cursorX
	f.state.selectionStart = snapshot.selectionStart
	f.state.selectionEnd = snapshot.selectionEnd
}

// hasSelection returns true if there is currently selected text
func (f *finder) hasSelection() bool {
	return f.state.selectionStart != f.state.selectionEnd
}

// clearSelection clears the current text selection
func (f *finder) clearSelection() {
	f.state.selectionStart = 0
	f.state.selectionEnd = 0
}

// deleteSelection removes the selected text and returns true if text was deleted
func (f *finder) deleteSelection() bool {
	if !f.hasSelection() {
		return false
	}

	start, end := f.getSelectionRange()
	f.state.input = append(f.state.input[:start], f.state.input[end:]...)

	// Update cursor position
	f.state.x = start
	f.state.cursorX = runewidth.StringWidth(string(f.state.input[:start]))

	f.clearSelection()
	return true
}

// getSelectionRange returns the start and end indices in the correct order
func (f *finder) getSelectionRange() (start, end int) {
	if f.state.selectionStart < f.state.selectionEnd {
		return f.state.selectionStart, f.state.selectionEnd
	}
	return f.state.selectionEnd, f.state.selectionStart
}

// getSelectedText returns the currently selected text
func (f *finder) getSelectedText() string {
	if !f.hasSelection() {
		return ""
	}
	start, end := f.getSelectionRange()
	return string(f.state.input[start:end])
}

// selectWord selects the word at the given position
func (f *finder) selectWord(pos int) {
	if pos < 0 || pos >= len(f.state.input) {
		return
	}

	// Find word boundaries
	start := pos
	end := pos

	// Move start backward to beginning of word
	for start > 0 && !unicode.IsSpace(f.state.input[start-1]) {
		start--
	}

	// Move end forward to end of word
	for end < len(f.state.input) && !unicode.IsSpace(f.state.input[end]) {
		end++
	}

	f.state.selectionStart = start
	f.state.selectionEnd = end
}

// runeIndexToScreenX converts a rune index in the input to screen X position
func (f *finder) runeIndexToScreenX(runeIdx int) int {
	if runeIdx <= 0 {
		return 0
	}
	if runeIdx > len(f.state.input) {
		runeIdx = len(f.state.input)
	}
	return runewidth.StringWidth(string(f.state.input[:runeIdx]))
}

// screenXToRuneIndex converts a screen X position to a rune index in the input
func (f *finder) screenXToRuneIndex(screenX int) int {
	if screenX <= 0 {
		return 0
	}

	currentWidth := 0
	for i, r := range f.state.input {
		rw := runewidth.RuneWidth(r)
		if currentWidth+rw > screenX {
			// Return the closest position
			if screenX-currentWidth < rw/2 {
				return i
			}
			return i + 1
		}
		currentWidth += rw
	}
	return len(f.state.input)
}

// pasteFromClipboard pastes text from clipboard at cursor position
func (f *finder) pasteFromClipboard() {
	text, err := clipboard.ReadAll()
	if err == nil && text != "" {
		f.saveUndoState()
		// Delete selection if any
		f.deleteSelection()

		// Insert clipboard text at cursor
		x := f.state.x
		newRunes := []rune(text)
		f.state.input = append(f.state.input[:x], append(newRunes, f.state.input[x:]...)...)
		f.state.x += len(newRunes)
		f.state.cursorX += runewidth.StringWidth(text)
	}
}

// cutToClipboard cuts selected text to clipboard
func (f *finder) cutToClipboard() {
	if f.hasSelection() {
		selectedText := f.getSelectedText()
		if err := clipboard.WriteAll(selectedText); err != nil {
			// Clipboard operation failed, but we don't abort the operation
			// Just skip the clipboard write and continue with the cut
		}
		f.saveUndoState()
		f.deleteSelection()
	}
}

// copyToClipboard copies selected text to clipboard
func (f *finder) copyToClipboard() {
	if f.hasSelection() {
		selectedText := f.getSelectedText()
		if err := clipboard.WriteAll(selectedText); err != nil {
			// Clipboard operation failed, silently ignore
			// User can retry if needed
		}
	}
}

// handleMouseEvent processes mouse events for text selection and clipboard operations
func (f *finder) handleMouseEvent(mouseX, mouseY int, buttons tcell.ButtonMask, layout Layout) {
	// Check if mouse is within the prompt area
	promptLinePad := len(f.opt.promptString)
	if mouseY != layout.prompt.y {
		// Click outside prompt line, ignore
		return
	}

	// Calculate relative X position in the input field
	relativeX := mouseX - layout.prompt.x - promptLinePad
	if relativeX < 0 {
		relativeX = 0
	}

	switch {
	case buttons&tcell.Button1 != 0:
		// Left mouse button pressed/held
		if !f.state.isDragging {
			// This is a mouse down event
			currentTime := time.Now()
			timeSinceLastClick := currentTime.Sub(f.state.mouseDownTime)

			// Check for double-click (within doubleClickThreshold)
			if timeSinceLastClick < doubleClickThreshold && f.state.mouseDownX == relativeX {
				// Double-click: select word
				f.state.lastClickCount = 2
				runeIdx := f.screenXToRuneIndex(relativeX)
				f.selectWord(runeIdx)
				f.state.mouseDownTime = time.Time{} // Reset to prevent triple-click
			} else {
				// Single click: start selection
				f.state.lastClickCount = 1
				f.state.mouseDownTime = currentTime
				f.state.mouseDownX = relativeX
				f.state.isDragging = true

				// Position cursor at click location
				runeIdx := f.screenXToRuneIndex(relativeX)
				f.state.x = runeIdx
				f.state.cursorX = f.runeIndexToScreenX(runeIdx)

				// Start selection
				f.state.selectionStart = runeIdx
				f.state.selectionEnd = runeIdx
			}
		} else {
			// Mouse is being dragged: update selection
			runeIdx := f.screenXToRuneIndex(relativeX)
			f.state.selectionEnd = runeIdx
			f.state.x = runeIdx
			f.state.cursorX = f.runeIndexToScreenX(runeIdx)
		}

	case buttons&tcell.Button3 != 0:
		// Right mouse button: paste from clipboard
		f.pasteFromClipboard()

	case buttons == tcell.ButtonNone:
		// Mouse button released
		if f.state.isDragging {
			f.state.isDragging = false
			// If selection start equals end, clear selection
			if f.state.selectionStart == f.state.selectionEnd {
				f.clearSelection()
			}
		}
	}
}
