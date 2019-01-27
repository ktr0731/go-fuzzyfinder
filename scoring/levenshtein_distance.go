// package scoring provides APIs that calculates similarity scores between two strings.
//
// Scoring consists of several parts as following.
//
//   - The distance between two strings by Levenshtein distance (edit distance) algorithm.
//   - Gap penalty by Smith-Waterman algorithm.
//
package scoring

import "fmt"

// Calculate calculates a similarity score between s1 and s2.
func Calculate(s1, s2 string) int {
	// return levenshteinDistance([]rune(s1), []rune(s2))
	// return needlemanWunsch([]rune(s1), []rune(s2))
	return smithWaterman([]rune(s1), []rune(s2))
}

func levenshteinDistance(s1, s2 []rune) int {
	m := make([][]int, len(s1)+1)
	for i := 0; i <= len(s1); i++ {
		m[i] = make([]int, len(s2)+1)
	}

	for i := 0; i <= len(s1); i++ {
		m[i][0] = i
	}
	for j := 0; j <= len(s2); j++ {
		m[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			ins := m[i-1][j] + 1
			del := m[i][j-1] + 1
			var rep int
			if s1[i-1] != s2[j-1] {
				rep = m[i-1][j-1] + 1
			} else {
				rep = m[i-1][j-1]
			}
			m[i][j] = min(ins, del, rep)
		}
	}

	fmt.Printf("%4c   ", '|')
	for i := 0; i < len(s2); i++ {
		fmt.Printf("%0c ", s2[i])
	}
	fmt.Printf("\n--------------------\n")

	fmt.Print("   | ")
	for i := 0; i <= len(s1); i++ {
		if i != 0 {
			fmt.Printf("%3c| ", s1[i-1])
		}
		for j := 0; j <= len(s2); j++ {
			fmt.Printf("%0d ", m[i][j])
		}
		fmt.Println()
	}
	return m[len(s1)][len(s2)]
}

// min returns the smallest number from passed args.
// If the number of args is 0, it always returns 0.
func min(n ...int) (min int) {
	if len(n) == 0 {
		return 0
	}
	min = n[0]
	for _, a := range n[1:] {
		if a < min {
			min = a
		}
	}
	return
}

// max returns the biggest number from passed args.
// If the number of args is 0, it always returns 0.
func max(n ...int) (min int) {
	if len(n) == 0 {
		return 0
	}
	min = n[0]
	for _, a := range n[1:] {
		if a > min {
			min = a
		}
	}
	return
}
