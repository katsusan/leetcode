package code

func countCompleteDayPairs(hours []int) int {
	m := make(map[int]int)
	pairs := 0

	for _, h := range hours {
		r := h % 24
		pairs += m[(24-r)%24] // must be %24 here to handle 24x case
		m[r]++
	}

	return pairs
}
