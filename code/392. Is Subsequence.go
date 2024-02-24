package code

import "strings"

func isSubsequence(s string, t string) bool {
	if len(s) == 0 {
		return true
	}

	chIdx := strings.IndexByte(t, s[0])
	if chIdx == -1 {
		return false
	}
	return isSubsequence(s[1:], t[chIdx+1:])
}

func isSubsequence2(s string, t string) bool {
	k := 0
	for i := 0; i < len(t) && k < len(s); i++ {
		if t[i] == s[k] {
			k++
		}
	}
	return k == len(s)
}
