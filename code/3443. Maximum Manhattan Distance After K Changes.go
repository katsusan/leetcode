package code

func maxDistance3443(s string, k int) int {
	mhDist := 0
	m := [91]int{} // stores counts of N,S,E,W
	for i := range s {
		m[s[i]]++
		minNS := min(m['N'], m['S'])
		minEW := min(m['E'], m['W'])
		if minNS+minEW <= k {
			// we can flip all moves we want
			mhDist = max(mhDist, i+1)
		} else {
			mhDist = max(mhDist, absDist(m['N'], m['S'])+absDist(m['E'], m['W'])+k*2)
		}
	}

	return mhDist
}

func absDist(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}
