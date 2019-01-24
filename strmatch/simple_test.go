package strmatch

import (
	"testing"
)

func Test_simpleMatch(t *testing.T) {
	in := "ink now"
	slice := []string{
		"WHITE ALBUM",
		"SOUND OF DESTINY",
		"Twinkle Snow",
	}

	simpleMatch(in, slice, option{})
}

func BenchmarkSimple(b *testing.B) {
	in := "ink"
	slice := []string{
		"WHITE ALBUM",
		"SOUND OF DESTINY",
		"Twinkle Snow",
	}
	for i := 0; i < b.N; i++ {
		simpleMatch(in, slice, option{})
	}
}
