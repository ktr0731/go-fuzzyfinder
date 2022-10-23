package scoring

import (
	"fmt"
	"os"
	"testing"
)

func Test_smithWaterman(t *testing.T) {
	t.Parallel()

	old := os.Getenv("DEBUG")
	os.Setenv("DEBUG", "true")
	defer os.Setenv("DEBUG", old)

	cases := []struct {
		s1, s2        string
		expectedScore int
		expectedPos   [2]int
	}{
		{"TACGGGCCCGCTA", "TAGCCCTA", 78, [2]int{0, 12}},
		{"TACGGG-CCCGCTA", "TAGCCCTA", 56, [2]int{0, 13}},
		{"FLY ME TO THE MOON", "MEON", 10, [2]int{4, 17}},
	}

	for _, c := range cases {
		c := c
		name := fmt.Sprintf("%s-%s", c.s1, c.s2)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			score, pos := smithWaterman([]rune(c.s1), []rune(c.s2))
			if score != c.expectedScore {
				t.Errorf("expected 78, but got %d", score)
			}
			if pos != c.expectedPos {
				t.Errorf("expected %v, but got %v", c.expectedPos, pos)
			}
		})
	}
}

func Benchmark_smithWaterman(b *testing.B) {
	for i := 0; i < b.N; i++ {
		smithWaterman([]rune("TACGGGCCCGCTA"), []rune("TAGCCCTA"))
	}
}
