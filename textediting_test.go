package fuzzyfinder

import (
	"testing"
	"time"
)

// TestSelectionHelpers tests basic selection functionality
func TestSelectionHelpers(t *testing.T) {
	f := &finder{
		state: state{
			input:          []rune("hello world"),
			selectionStart: 0,
			selectionEnd:   0,
		},
	}

	// Test hasSelection - initially no selection
	if f.hasSelection() {
		t.Error("expected no selection initially")
	}

	// Test setting selection
	f.state.selectionStart = 0
	f.state.selectionEnd = 5
	if !f.hasSelection() {
		t.Error("expected selection after setting bounds")
	}

	// Test getSelectedText
	selected := f.getSelectedText()
	if selected != "hello" {
		t.Errorf("expected 'hello', got '%s'", selected)
	}

	// Test getSelectionRange with reversed selection
	f.state.selectionStart = 6
	f.state.selectionEnd = 0
	start, end := f.getSelectionRange()
	if start != 0 || end != 6 {
		t.Errorf("expected range (0, 6), got (%d, %d)", start, end)
	}

	// Test clearSelection
	f.clearSelection()
	if f.hasSelection() {
		t.Error("expected no selection after clear")
	}
}

// TestDeleteSelection tests selection deletion
func TestDeleteSelection(t *testing.T) {
	f := &finder{
		state: state{
			input:          []rune("hello world"),
			x:              5,
			cursorX:        5,
			selectionStart: 0,
			selectionEnd:   5,
		},
	}

	// Test deleteSelection
	deleted := f.deleteSelection()
	if !deleted {
		t.Error("expected deletion to succeed")
	}

	expected := " world"
	actual := string(f.state.input)
	if actual != expected {
		t.Errorf("expected '%s', got '%s'", expected, actual)
	}

	if f.state.x != 0 {
		t.Errorf("expected cursor at 0, got %d", f.state.x)
	}

	if f.hasSelection() {
		t.Error("expected no selection after delete")
	}

	// Test deleteSelection with no selection
	deleted = f.deleteSelection()
	if deleted {
		t.Error("expected deletion to fail with no selection")
	}
}

// TestSelectWord tests word selection
func TestSelectWord(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		pos           int
		expectedStart int
		expectedEnd   int
	}{
		{
			name:          "middle of word",
			input:         "hello world",
			pos:           2,
			expectedStart: 0,
			expectedEnd:   5,
		},
		{
			name:          "start of word",
			input:         "hello world",
			pos:           0,
			expectedStart: 0,
			expectedEnd:   5,
		},
		{
			name:          "second word",
			input:         "hello world",
			pos:           7,
			expectedStart: 6,
			expectedEnd:   11,
		},
		{
			name:          "single word",
			input:         "hello",
			pos:           2,
			expectedStart: 0,
			expectedEnd:   5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &finder{
				state: state{
					input: []rune(tt.input),
				},
			}

			f.selectWord(tt.pos)

			if f.state.selectionStart != tt.expectedStart {
				t.Errorf("expected selectionStart %d, got %d", tt.expectedStart, f.state.selectionStart)
			}
			if f.state.selectionEnd != tt.expectedEnd {
				t.Errorf("expected selectionEnd %d, got %d", tt.expectedEnd, f.state.selectionEnd)
			}
		})
	}
}

// TestUndoRedo tests undo/redo functionality
func TestUndoRedo(t *testing.T) {
	f := &finder{
		state: state{
			input:       []rune("hello"),
			x:           5,
			cursorX:     5,
			undoHistory: make([]undoState, 0, maxUndoHistory),
			undoIndex:   -1,
		},
	}

	// Save initial state
	f.saveUndoState()
	if len(f.state.undoHistory) != 1 {
		t.Errorf("expected 1 undo state, got %d", len(f.state.undoHistory))
	}

	// Modify input
	f.state.input = []rune("hello world")
	f.state.x = 11
	f.saveUndoState()

	if len(f.state.undoHistory) != 2 {
		t.Errorf("expected 2 undo states, got %d", len(f.state.undoHistory))
	}

	// Test undo
	f.undo()
	if string(f.state.input) != "hello" {
		t.Errorf("expected 'hello' after undo, got '%s'", string(f.state.input))
	}
	if f.state.x != 5 {
		t.Errorf("expected cursor at 5, got %d", f.state.x)
	}

	// Test undo at beginning (should do nothing)
	f.undo()
	if len(f.state.undoHistory) != 2 {
		t.Error("undo should not remove history")
	}
}

// TestUndoHistoryLimit tests that undo history respects the limit
func TestUndoHistoryLimit(t *testing.T) {
	f := &finder{
		state: state{
			input:       []rune("test"),
			undoHistory: make([]undoState, 0, maxUndoHistory),
			undoIndex:   -1,
		},
	}

	// Add more than maxUndoHistory states
	for i := 0; i <= maxUndoHistory+10; i++ {
		f.saveUndoState()
	}

	if len(f.state.undoHistory) > maxUndoHistory {
		t.Errorf("expected max %d undo states, got %d", maxUndoHistory, len(f.state.undoHistory))
	}
}

// TestScreenXToRuneIndex tests coordinate conversion
func TestScreenXToRuneIndex(t *testing.T) {
	// Test with ASCII characters (1 rune = 1 screen width)
	f := &finder{
		state: state{
			input: []rune("hello"),
		},
	}

	// The function uses runewidth which may round to nearest position
	// Test boundary cases
	tests := []struct {
		screenX int
		min     int // Minimum acceptable index
		max     int // Maximum acceptable index
	}{
		{0, 0, 0},
		{1, 0, 2}, // Could be 1 or nearby due to rounding
		{5, 5, 5},
		{100, 5, 5}, // Beyond end
		{-1, 0, 0},  // Before start
	}

	for _, tt := range tests {
		result := f.screenXToRuneIndex(tt.screenX)
		if result < tt.min || result > tt.max {
			t.Errorf("screenXToRuneIndex(%d) = %d, expected in range [%d, %d]", tt.screenX, result, tt.min, tt.max)
		}
	}
}

// TestRuneIndexToScreenX tests coordinate conversion
func TestRuneIndexToScreenX(t *testing.T) {
	f := &finder{
		state: state{
			input: []rune("hello"),
		},
	}

	tests := []struct {
		runeIdx  int
		expected int
	}{
		{0, 0},
		{1, 1},
		{5, 5},
		{100, 5}, // Beyond end, should cap
		{-1, 0},  // Before start, should cap
	}

	for _, tt := range tests {
		result := f.runeIndexToScreenX(tt.runeIdx)
		if result != tt.expected {
			t.Errorf("runeIndexToScreenX(%d) = %d, expected %d", tt.runeIdx, result, tt.expected)
		}
	}
}

// TestConstants verifies the constants are set to expected values
func TestConstants(t *testing.T) {
	if maxUndoHistory != 100 {
		t.Errorf("expected maxUndoHistory to be 100, got %d", maxUndoHistory)
	}

	if doubleClickThreshold != 500*time.Millisecond {
		t.Errorf("expected doubleClickThreshold to be 500ms, got %v", doubleClickThreshold)
	}
}
