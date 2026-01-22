package fuzzyfinder

import "testing"

func TestRect_SplitVertical(t *testing.T) {
	tests := []struct {
		name      string
		r         rect
		ratio     float64
		wantLeft  rect
		wantRight rect
	}{
		{
			name:      "split 50/50",
			r:         rect{x: 0, y: 0, width: 100, height: 20},
			ratio:     0.5,
			wantLeft:  rect{x: 0, y: 0, width: 50, height: 20},
			wantRight: rect{x: 50, y: 0, width: 50, height: 20},
		},
		{
			name:      "split 30/70",
			r:         rect{x: 10, y: 5, width: 100, height: 20},
			ratio:     0.3,
			wantLeft:  rect{x: 10, y: 5, width: 30, height: 20},
			wantRight: rect{x: 40, y: 5, width: 70, height: 20},
		},
		{
			name:      "invalid ratio (too low) defaults to 50/50",
			r:         rect{x: 0, y: 0, width: 100, height: 20},
			ratio:     0.0,
			wantLeft:  rect{x: 0, y: 0, width: 50, height: 20},
			wantRight: rect{x: 50, y: 0, width: 50, height: 20},
		},
		{
			name:      "invalid ratio (too high) defaults to 50/50",
			r:         rect{x: 0, y: 0, width: 100, height: 20},
			ratio:     1.0,
			wantLeft:  rect{x: 0, y: 0, width: 50, height: 20},
			wantRight: rect{x: 50, y: 0, width: 50, height: 20},
		},
		{
			name:      "odd width",
			r:         rect{x: 0, y: 0, width: 99, height: 20},
			ratio:     0.5,
			wantLeft:  rect{x: 0, y: 0, width: 49, height: 20},
			wantRight: rect{x: 49, y: 0, width: 50, height: 20},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLeft, gotRight := tt.r.splitVertical(tt.ratio)
			if gotLeft != tt.wantLeft {
				t.Errorf("splitVertical() left = %v, want %v", gotLeft, tt.wantLeft)
			}
			if gotRight != tt.wantRight {
				t.Errorf("splitVertical() right = %v, want %v", gotRight, tt.wantRight)
			}
		})
	}
}

func TestRect_Inset(t *testing.T) {
	tests := []struct {
		name   string
		r      rect
		margin int
		want   rect
	}{
		{
			name:   "inset by 1",
			r:      rect{x: 0, y: 0, width: 100, height: 50},
			margin: 1,
			want:   rect{x: 1, y: 1, width: 98, height: 48},
		},
		{
			name:   "inset by 2",
			r:      rect{x: 10, y: 5, width: 100, height: 50},
			margin: 2,
			want:   rect{x: 12, y: 7, width: 96, height: 46},
		},
		{
			name:   "inset by 0",
			r:      rect{x: 10, y: 5, width: 100, height: 50},
			margin: 0,
			want:   rect{x: 10, y: 5, width: 100, height: 50},
		},
		{
			name:   "inset large margin",
			r:      rect{x: 0, y: 0, width: 10, height: 10},
			margin: 5,
			want:   rect{x: 5, y: 5, width: 0, height: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.r.inset(tt.margin)
			if got != tt.want {
				t.Errorf("inset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRect_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		r    rect
		want bool
	}{
		{
			name: "normal rectangle",
			r:    rect{x: 0, y: 0, width: 100, height: 50},
			want: false,
		},
		{
			name: "zero width",
			r:    rect{x: 0, y: 0, width: 0, height: 50},
			want: true,
		},
		{
			name: "zero height",
			r:    rect{x: 0, y: 0, width: 100, height: 0},
			want: true,
		},
		{
			name: "negative width",
			r:    rect{x: 0, y: 0, width: -10, height: 50},
			want: true,
		},
		{
			name: "negative height",
			r:    rect{x: 0, y: 0, width: 100, height: -5},
			want: true,
		},
		{
			name: "1x1 rectangle",
			r:    rect{x: 0, y: 0, width: 1, height: 1},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.r.isEmpty()
			if got != tt.want {
				t.Errorf("isEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLayout_Validate(t *testing.T) {
	tests := []struct {
		name    string
		layout  Layout
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid layout",
			layout: Layout{
				content:      rect{x: 0, y: 0, width: 100, height: 50},
				innerContent: rect{x: 1, y: 1, width: 98, height: 48},
				list:         rect{x: 1, y: 1, width: 48, height: 48},
				items:        rect{x: 1, y: 1, width: 48, height: 45},
				hasPreview:   true,
				preview:      rect{x: 50, y: 1, width: 48, height: 48},
			},
			wantErr: false,
		},
		{
			name: "empty content area",
			layout: Layout{
				content:      rect{x: 0, y: 0, width: 0, height: 0},
				innerContent: rect{x: 1, y: 1, width: 98, height: 48},
				list:         rect{x: 1, y: 1, width: 48, height: 48},
				items:        rect{x: 1, y: 1, width: 48, height: 45},
			},
			wantErr: true,
			errMsg:  "content area is too small",
		},
		{
			name: "empty inner content",
			layout: Layout{
				content:      rect{x: 0, y: 0, width: 100, height: 50},
				innerContent: rect{x: 0, y: 0, width: 0, height: 0},
				list:         rect{x: 1, y: 1, width: 48, height: 48},
				items:        rect{x: 1, y: 1, width: 48, height: 45},
			},
			wantErr: true,
			errMsg:  "inner content area is too small (border may be too large)",
		},
		{
			name: "empty list area",
			layout: Layout{
				content:      rect{x: 0, y: 0, width: 100, height: 50},
				innerContent: rect{x: 1, y: 1, width: 98, height: 48},
				list:         rect{x: 0, y: 0, width: 0, height: 0},
				items:        rect{x: 1, y: 1, width: 48, height: 45},
			},
			wantErr: true,
			errMsg:  "list area is too small",
		},
		{
			name: "empty items area",
			layout: Layout{
				content:      rect{x: 0, y: 0, width: 100, height: 50},
				innerContent: rect{x: 1, y: 1, width: 98, height: 48},
				list:         rect{x: 1, y: 1, width: 48, height: 48},
				items:        rect{x: 0, y: 0, width: 0, height: 0},
			},
			wantErr: true,
			errMsg:  "items area is too small (need at least 1 line for items)",
		},
		{
			name: "list too narrow",
			layout: Layout{
				content:      rect{x: 0, y: 0, width: 100, height: 50},
				innerContent: rect{x: 1, y: 1, width: 98, height: 48},
				list:         rect{x: 1, y: 1, width: 5, height: 48},
				items:        rect{x: 1, y: 1, width: 5, height: 45},
			},
			wantErr: true,
			errMsg:  "terminal is too narrow (need at least 10 columns for list)",
		},
		{
			name: "preview enabled but empty",
			layout: Layout{
				content:      rect{x: 0, y: 0, width: 100, height: 50},
				innerContent: rect{x: 1, y: 1, width: 98, height: 48},
				list:         rect{x: 1, y: 1, width: 48, height: 48},
				items:        rect{x: 1, y: 1, width: 48, height: 45},
				hasPreview:   true,
				preview:      rect{x: 0, y: 0, width: 0, height: 0},
			},
			wantErr: true,
			errMsg:  "preview area is too small",
		},
		{
			name: "preview too narrow",
			layout: Layout{
				content:      rect{x: 0, y: 0, width: 100, height: 50},
				innerContent: rect{x: 1, y: 1, width: 98, height: 48},
				list:         rect{x: 1, y: 1, width: 48, height: 48},
				items:        rect{x: 1, y: 1, width: 48, height: 45},
				hasPreview:   true,
				preview:      rect{x: 50, y: 1, width: 5, height: 48},
			},
			wantErr: true,
			errMsg:  "preview area is too narrow (need at least 10 columns)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.layout.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestFinder_ComputeLayout(t *testing.T) {
	tests := []struct {
		name         string
		termWidth    int
		termHeight   int
		optWidth     int
		optHeight    int
		optBorder    bool
		optHeader    string
		optPreview   bool
		wantErr      bool
		validateFunc func(t *testing.T, layout Layout)
	}{
		{
			name:       "basic layout without options",
			termWidth:  80,
			termHeight: 24,
			optBorder:  false,
			optPreview: false,
			wantErr:    false,
			validateFunc: func(t *testing.T, layout Layout) {
				if layout.terminal.width != 80 || layout.terminal.height != 24 {
					t.Errorf("terminal size = %dx%d, want 80x24", layout.terminal.width, layout.terminal.height)
				}
				if layout.content != layout.innerContent {
					t.Error("content and innerContent should be same without border")
				}
				if layout.list != layout.innerContent {
					t.Error("list should equal innerContent without preview")
				}
				if !layout.preview.isEmpty() {
					t.Error("preview should be empty when not enabled")
				}
			},
		},
		{
			name:       "layout with border",
			termWidth:  80,
			termHeight: 24,
			optBorder:  true,
			optPreview: false,
			wantErr:    false,
			validateFunc: func(t *testing.T, layout Layout) {
				if layout.innerContent.width != layout.content.width-2 {
					t.Errorf("innerContent width = %d, want %d", layout.innerContent.width, layout.content.width-2)
				}
				if layout.innerContent.height != layout.content.height-2 {
					t.Errorf("innerContent height = %d, want %d", layout.innerContent.height, layout.content.height-2)
				}
			},
		},
		{
			name:       "layout with preview",
			termWidth:  80,
			termHeight: 24,
			optBorder:  false,
			optPreview: true,
			wantErr:    false,
			validateFunc: func(t *testing.T, layout Layout) {
				if layout.list.width+layout.preview.width != layout.innerContent.width {
					t.Errorf("list(%d) + preview(%d) != innerContent(%d)",
						layout.list.width, layout.preview.width, layout.innerContent.width)
				}
				if layout.preview.isEmpty() {
					t.Error("preview should not be empty when enabled")
				}
			},
		},
		{
			name:       "layout with constrained height",
			termWidth:  80,
			termHeight: 24,
			optHeight:  10,
			optBorder:  false,
			optPreview: false,
			wantErr:    false,
			validateFunc: func(t *testing.T, layout Layout) {
				if layout.content.height != 10 {
					t.Errorf("content height = %d, want 10", layout.content.height)
				}
				// When height is constrained, content should be positioned at bottom
				expectedY := 24 - 10 // termHeight - contentHeight
				if layout.content.y != expectedY {
					t.Errorf("content y position = %d, want %d (bottom-aligned)", layout.content.y, expectedY)
				}
			},
		},
		{
			name:       "layout with constrained width",
			termWidth:  80,
			termHeight: 24,
			optWidth:   40,
			optBorder:  false,
			optPreview: false,
			wantErr:    false,
			validateFunc: func(t *testing.T, layout Layout) {
				if layout.content.width != 40 {
					t.Errorf("content width = %d, want 40", layout.content.width)
				}
			},
		},
		{
			name:       "layout with header",
			termWidth:  80,
			termHeight: 24,
			optBorder:  false,
			optHeader:  "Test Header",
			optPreview: false,
			wantErr:    false,
			validateFunc: func(t *testing.T, layout Layout) {
				if layout.header.isEmpty() {
					t.Error("header should not be empty when set")
				}
				if layout.header.height != 1 {
					t.Errorf("header height = %d, want 1", layout.header.height)
				}
			},
		},
		{
			name:         "layout too small",
			termWidth:    5,
			termHeight:   2,
			optBorder:    false,
			optPreview:   false,
			wantErr:      true,
			validateFunc: func(t *testing.T, layout Layout) {},
		},
		{
			name:       "combined: border + preview + header + constraints",
			termWidth:  100,
			termHeight: 30,
			optWidth:   80,
			optHeight:  20,
			optBorder:  true,
			optHeader:  "Header",
			optPreview: true,
			wantErr:    false,
			validateFunc: func(t *testing.T, layout Layout) {
				if layout.content.width != 80 || layout.content.height != 20 {
					t.Errorf("content = %dx%d, want 80x20", layout.content.width, layout.content.height)
				}
				if !layout.hasBorder || !layout.hasPreview || !layout.hasHeader {
					t.Error("all features should be enabled")
				}
				// Verify all areas are positioned correctly
				if layout.items.isEmpty() {
					t.Error("items area should not be empty")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFinder()
			term := f.UseMockedTerminalV2()
			term.SetSize(tt.termWidth, tt.termHeight)

			f.opt = &opt{
				width:  tt.optWidth,
				height: tt.optHeight,
				border: tt.optBorder,
				header: tt.optHeader,
			}

			if tt.optPreview {
				f.opt.previewFunc = func(i, w, h int) string { return "preview" }
			}

			layout, err := f.computeLayout()
			if (err != nil) != tt.wantErr {
				t.Errorf("computeLayout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Validate the layout
				if err := layout.validate(); err != nil {
					t.Errorf("layout validation failed: %v", err)
				}

				// Run custom validation
				tt.validateFunc(t, layout)
			}
		})
	}
}
