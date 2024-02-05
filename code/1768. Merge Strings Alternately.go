package code

/*
You are given two strings word1 and word2. Merge the strings by adding letters in alternating order, starting with word1.
If a string is longer than the other, append the additional letters onto the end of the merged string.

TLDR,合并word1和word2:)
*/

func mergeAlternately(word1 string, word2 string) string {
	var b = make([]byte, len(word1)+len(word2))
	var idx int
	for idx = 0; idx < min(len(word1), len(word2)); idx++ {
		b[2*idx], b[2*idx+1] = word1[idx], word2[idx]
	}

	if len(word1) < len(word2) {
		copy(b[2*idx:], []byte(word2[idx:]))
	} else {
		copy(b[2*idx:], []byte(word1[idx:]))
	}
	return string(b)
}
