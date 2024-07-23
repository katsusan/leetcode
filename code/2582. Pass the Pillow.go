package code

func passThePillow(n int, time int) int {
	// for n = 4, we can think a infinite sequence: 1 2 3 4 3 2 1 2 3 4...
	// thus for every 6 times we come to the starting point
	r := time % (2 * (n - 1))
	if r < n {
		return r + 1
	}

	return 2*n - 1 - r
}
