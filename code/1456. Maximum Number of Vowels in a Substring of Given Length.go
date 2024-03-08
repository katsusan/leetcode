package code

import "strings"

func maxVowels(s string, k int) int {
	if len(s) < k {
		return CountVowels(s)
	}

	wd := CountVowels(s[0:k])
	maxCnt := wd
	for i := 1; i <= len(s)-k; i++ {
		if IsVowel(s[i-1]) {
			wd--
		}

		if IsVowel(s[i-1+k]) {
			wd++
		}

		maxCnt = max(maxCnt, wd)
	}

	return maxCnt
}

func CountVowels(s string) int {
	cnt := 0
	for i := range s {
		if IsVowel(s[i]) {
			cnt++
		}
	}
	return cnt
}

func IsVowel(b byte) bool {
	return strings.IndexByte("aeiou", b) != -1
}
