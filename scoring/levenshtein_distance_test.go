package scoring

import (
	"testing"
)

func TestLevenshteinDistance(t *testing.T) {
	s1 := []rune("kitten")
	s2 := []rune("sitting")

	const expected = 29
	score := levenshteinDistance(s1, s2)
	if score != expected {
		t.Errorf("expected score: %d, but got %d", expected, score)
	}
}
