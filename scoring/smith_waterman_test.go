package scoring

import (
	"testing"
)

func TestSW(t *testing.T) {
	score := smithWaterman([]rune("GGTTGACTA"), []rune("TGTTACGG"))
	if score != 13 {
		t.Errorf("expected 13, but got %d", score)
	}
}
