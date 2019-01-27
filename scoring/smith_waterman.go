package scoring

import (
	"flag"
	"fmt"
)

func smithWaterman(s1, s2 []rune) int {
	// Penalties.
	const (
		match      = 3
		mismatch   = 3
		gap        = 2
		openingGap = 10
	)

	m := make([][]int, len(s1)+1)
	for i := 0; i <= len(s1); i++ {
		m[i] = make([]int, len(s2)+1)
	}

	// TODO: 初期状態がおかしい

	var maxScore int
	gapOpening := true
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			s1gap := m[i-1][j] - gap
			s2gap := m[i][j-1] - gap
			if gapOpening {
				s1gap -= openingGap
				s2gap -= openingGap
			}
			var matchScore int
			if s1[i-1] != s2[j-1] {
				matchScore = m[i-1][j-1] - match
			} else {
				matchScore = m[i-1][j-1] + mismatch
			}
			m[i][j] = max(s1gap, s2gap, matchScore, 0)
			maxScore = max(m[i][j], maxScore)
			if m[i][j] == matchScore {
				gapOpening = true
			} else {
				gapOpening = false
			}
		}
	}

	if isInTesting() {
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

func isInTesting() bool {
	return flag.Lookup("test.v") != nil
}
