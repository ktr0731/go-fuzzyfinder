package scoring

import (
	"fmt"
	"testing"
)

func Test_smithWaterman(t *testing.T) {
	cases := []struct {
		s1, s2   string
		expected int
	}{
		{"TACGGGCCCGCTA", "TAGCCCTA", 78},
		{"TACGGG-CCCGCTA", "TAGCCCTA", 56},
	}

	for _, c := range cases {
		c := c
		name := fmt.Sprintf("%s-%s", c.s1, c.s2)
		t.Run(name, func(t *testing.T) {
			score := smithWaterman([]rune(c.s1), []rune(c.s2))
			if score != c.expected {
				t.Errorf("expected 78, but got %d", score)
			}
		})
	}
}

func Benchmark_smithWaterman(b *testing.B) {
	for i := 0; i < b.N; i++ {
		smithWaterman([]rune("TACGGGCCCGCTA"), []rune("TAGCCCTA"))
	}
}
