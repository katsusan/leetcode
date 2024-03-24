package code

import (
	"maps"
	"slices"
)

func closeStrings(word1 string, word2 string) bool {
	if len(word1) != len(word2) {
		return false
	}

	freqM1 := freqMap(word1)
	freqM2 := freqMap(word2)

	charSet1 := make(map[byte]struct{}, 26)
	freqSet1 := make([]int, 0, 26)

	charSet2 := make(map[byte]struct{}, 26)
	freqSet2 := make([]int, 0, 26)

	for c, freq := range freqM1 {
		charSet1[c] = struct{}{}
		freqSet1 = append(freqSet1, freq)
	}

	for c, freq := range freqM2 {
		charSet2[c] = struct{}{}
		freqSet2 = append(freqSet2, freq)
	}

	slices.Sort(freqSet1)
	slices.Sort(freqSet2)

	return maps.Equal(charSet1, charSet2) && slices.Equal(freqSet1, freqSet2)
}

func freqMap(word string) map[byte]int {
	s := make(map[byte]int, 26)
	for _, b := range []byte(word) {
		s[b]++
	}
	return s
}
