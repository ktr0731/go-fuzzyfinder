package scoring

import (
	"fmt"
	"os"
)

func smithWaterman(s1, s2 []rune) int {
	const (
		openGap = 5 // Gap opening penalty.
		extGap  = 1 // Gap extension penalty.

		matchScore    = 1
		mismatchScore = 1

		firstCharBonus = 1 // The first char of s1 is equal to s2's one.
	)

	// The scoring matrix.
	m := make([][]int, len(s1)+1)
	for i := 0; i <= len(s1); i++ {
		m[i] = make([]int, len(s2)+1)
	}

	if s1[0] == s2[0] {
		m[1][1] = firstCharBonus
	}

	var maxScore int
	gapOpening := true
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			// TODO: len(s) >= len(in) なので対象文字列にギャップは発生しない
			s1gap, s2gap := m[i-1][j], m[i][j-1]
			if gapOpening {
				s1gap -= openGap
				s2gap -= openGap
			} else {
				s1gap -= extGap
				s2gap -= extGap
			}
			var score int
			if s1[i-1] != s2[j-1] {
				score = m[i-1][j-1] - matchScore
			} else {
				score = m[i-1][j-1] + mismatchScore
			}
			m[i][j] += max(s1gap, s2gap, score, 0)

			// Update the max score.
			maxScore = max(m[i][j], maxScore)

			gapAdded := m[i][j] != score
			if gapAdded {
				gapOpening = false
			} else {
				gapOpening = true
			}
		}
	}

	if isDebug() {
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
	}

	return maxScore
}

func isDebug() bool {
	return os.Getenv("DEBUG") != ""
}
