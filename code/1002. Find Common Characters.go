package code

func commonChars(words []string) []string {
	charCnt := make([]int, 26)
	for _, c := range words[0] {
		charCnt[c-'a']++
	}
	for i := 1; i < len(words); i++ {
		tmpM := make([]int, 26)
		for _, c := range words[i] {
			tmpM[c-'a']++
		}
		charCnt = intersect(charCnt, tmpM)
	}

	cm := []string{}
	for i := range charCnt {
		s := string('a' + byte(i))
		n := charCnt[i]
		for j := 0; j < n; j++ {
			cm = append(cm, s)
		}
	}
	return cm
}

func intersect(m1, m2 []int) []int {
	for i := range m1 {
		m1[i] = min(m1[i], m2[i])
	}
	return m1
}
