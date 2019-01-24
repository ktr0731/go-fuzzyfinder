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

var in = "ink"
var slice = []string{
	"WHITE ALBUM",
	"SOUND OF DESTINY",
	"Twinkle Snow",
}

func BenchmarkSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simpleMatch(in, slice, option{})
	}
}

func BenchmarkSimple2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simpleMatch2(in, slice, option{})
	}
}
