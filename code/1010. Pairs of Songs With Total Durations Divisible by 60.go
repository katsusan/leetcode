package code

func numPairsDivisibleBy60(time []int) int {
	r := make([]int, 60)
	pairs := 0

	for _, t := range time {
		x := t % 60
		pairs += r[(60-x)%60]
		r[x]++
	}

	return pairs
}
