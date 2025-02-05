package code

func areAlmostEqual(s1 string, s2 string) bool {
	n := len(s1)
	swap := [2]int{}
	swapidx := -1
	for i := 0; i < n; i++ {
		if s1[i] != s2[i] {
			swapidx++
			if swapidx > 1 {
				return false
			}
			swap[swapidx] = i
		}
	}

	return s1[swap[0]] == s2[swap[1]] && s1[swap[1]] == s2[swap[0]]
}
