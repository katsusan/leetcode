package code

func firstUniqChar(s string) int {
	var m = make(map[byte]int, 26)

	for i := range s {
		if _, ok := m[s[i]]; !ok {
			m[s[i]] = i
			continue
		} else {
			m[s[i]] = len(s)
		}
	}

	minIdx := len(s)
	for _, v := range m {
		if v < minIdx {
			minIdx = v
		}
	}

	if minIdx == len(s) {
		return -1
	}

	return minIdx
}
