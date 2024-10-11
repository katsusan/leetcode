package code

func checkInclusion(s1 string, s2 string) bool {
	if len(s1) > len(s2) {
		return false
	}

	wndCnt := make([]int, 26)
	for _, c := range s1 {
		wndCnt[int(c-'a')]--
	}

	wndsize := len(s1)
	for i := 0; i < wndsize; i++ {
		wndCnt[int(s2[i]-'a')]++
	}

	if isZero(wndCnt) {
		return true
	}

	for j := wndsize; j < len(s2); j++ {
		wndCnt[int(s2[j-wndsize]-'a')]--
		wndCnt[int(s2[j]-'a')]++
		if isZero(wndCnt) {
			return true
		}
	}

	return false
}

func isZero[T int](s []T) bool {
	for i := range s {
		if s[i] != 0 {
			return false
		}
	}

	return true
}
