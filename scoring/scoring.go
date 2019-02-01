// package scoring provides APIs that calculates similarity scores between two strings.
package scoring

// Calculate calculates a similarity score between s1 and s2.
func Calculate(s1, s2 string) int {
	return smithWaterman([]rune(s1), []rune(s2))
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
