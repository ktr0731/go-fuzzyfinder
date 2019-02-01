package scoring

import (
	"testing"
)

func Test_smithWaterman(t *testing.T) {
	t.Skip("TODO")
	score := smithWaterman([]rune("TACGGGCCCGCTA"), []rune("TAGCCCTA"))
	if score != 13 {
		t.Errorf("expected 13, but got %d", score)
	}

	// Expected align:
	//
	// TACGGGCCCGCTA
	// TA---GCC--CTA
	//
}

func Benchmark_smithWaterman(b *testing.B) {
	for i := 0; i < b.N; i++ {
		smithWaterman([]rune("TACGGGCCCGCTA"), []rune("TAGCCCTA"))
	}
}
