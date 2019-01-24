package strmatch

import "fmt"

func levenshteinDistanceMatch(in string, slice []string, opt option) (res []Matched) {
	panic("TODO")
	return
}

func levenshteinDistance(s1, s2 string) int {
	m := make([][]int, len(s1)+1)
	for i := 0; i <= len(s1); i++ {
		m[i] = make([]int, len(s2)+1)
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			insertCost := m[i-1][j] + 1
			deleteCost := m[i][j-1] + 1
			var replaceCost int
			if s1[i-1] != s2[j-1] {
				replaceCost = m[i-1][j-1] + 1
			} else {
				replaceCost = m[i-1][j-1]
			}

			min := insertCost
			if deleteCost < min {
				min = deleteCost
			}
			if replaceCost < min {
				min = replaceCost
			}
			m[i][j] = min
		}
	}

	fmt.Print("  ")
	for i := 0; i < len(s2); i++ {
		fmt.Printf("%c ", s2[i])
	}
	fmt.Println("\n-------------")
	for i := 1; i <= len(s1); i++ {
		fmt.Printf("%c ", s1[i-1])
		for j := 1; j <= len(s2); j++ {
			fmt.Print(m[i][j], " ")
		}
		fmt.Println()
	}
	return m[len(s1)][len(s2)]
}
