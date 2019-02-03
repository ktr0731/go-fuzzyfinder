package scoring

import (
	"fmt"
	"os"
	"unicode"
)

// smithWaterman calculates a simularity score between s1 and s2
// by smith-waterman algorithm. smith-waterman algorithm is one of
// local alignment algorithms and it uses dynamic programming.
//
// In this smith-waterman algorithm, we use the affine gap penalty.
// Please see https://en.wikipedia.org/wiki/Gap_penalty#Affine for additional
// information about the affine gap penalty.
//
// We calculate the gap penalty by the Gotoh's algorithm, which optimizes
// the calculation from O(M^2N) to O(MN).
// Please see ftp://150.128.97.71/pub/Bioinformatica/gotoh1982.pdf for more details.
func smithWaterman(s1, s2 []rune) int {
	const (
		openGap = 5 // Gap opening penalty.
		extGap  = 1 // Gap extension penalty.

		matchScore    = 5
		mismatchScore = 1

		firstCharBonus = 3 // The first char of s1 is equal to s2's one.
	)

	// The scoring matrix.
	m := make([][]int, len(s1)+1)
	// A matrix that calculates gap penalties until each position (i, j).
	P := make([][]int, len(s1)+1)
	for i := 0; i <= len(s1); i++ {
		m[i] = make([]int, len(s2)+1)
		P[i] = make([]int, len(s2)+1)
	}

	for i := 0; i <= len(s1); i++ {
		P[i][0] = -openGap - i*extGap
	}

	///

	bonus := make([]int, len(s1))
	bonus[0] = firstCharBonus
	prevCh := s1[0]
	prevIsDelimiter := isDelimiter(prevCh)
	for i, r := range s1[1:] {
		isDelimiter := isDelimiter(r)
		if prevIsDelimiter && !isDelimiter {
			bonus[i] = firstCharBonus
		}
		prevIsDelimiter = isDelimiter
	}

	///

	var maxScore int
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			p := P[i-1][j] + bonus[i-1]
			var score int
			if s1[i-1] != s2[j-1] {
				score = m[i-1][j-1] - mismatchScore
			} else {
				score = m[i-1][j-1] + matchScore + bonus[i-1]
			}
			m[i][j] += max(p, score, 0)

			P[i][j] = max(m[i-1][j]-openGap, P[i-1][j]-extGap)

			// Update the max score.
			maxScore = max(m[i][j], maxScore)
		}
	}

	if isDebug() {
		printSlice := func(m [][]int) {
			fmt.Printf("%4c     ", '|')
			for i := 0; i < len(s2); i++ {
				fmt.Printf("%3c ", s2[i])
			}
			fmt.Printf("\n-------------------------\n")

			fmt.Print("   | ")
			for i := 0; i <= len(s1); i++ {
				if i != 0 {
					fmt.Printf("%3c| ", s1[i-1])
				}
				for j := 0; j <= len(s2); j++ {
					fmt.Printf("%3d ", m[i][j])
				}
				fmt.Println()
			}
			println()
		}
		printSlice(m)
		printSlice(P)
	}

	// We adjust scores by the weight per one rune.
	return int(float32(maxScore) * (float32(maxScore) / float32(len(s1))))
}

func isDebug() bool {
	return os.Getenv("DEBUG") != ""
}

var delimiterRunes = map[rune]interface{}{
	'(': nil,
	'[': nil,
	'{': nil,
	'/': nil,
	'-': nil,
	'_': nil,
	'.': nil,
}

func isDelimiter(r rune) bool {
	if _, ok := delimiterRunes[r]; ok {
		return true
	}
	return unicode.IsSpace(r)
}
