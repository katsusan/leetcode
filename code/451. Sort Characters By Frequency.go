package code

/*
Given a string s, sort it in decreasing order based on the frequency of the characters. The frequency of a character is the number of times it appears in the string.

Return the sorted string. If there are multiple answers, return any of them.
*/

import (
	"bytes"
	"slices"
)

func frequencySort(s string) string {
	freqCnt := make(map[byte]int)
	for i := range s {
		if cnt, ok := freqCnt[s[i]]; ok {
			freqCnt[s[i]] = cnt + 1
		} else {
			freqCnt[s[i]] = 1
		}
	}

	type charBucket struct {
		elem byte
		cnt  int
	}

	var chars = make([]charBucket, len(freqCnt))
	var i int
	for e, c := range freqCnt {
		chars[i].elem = e
		chars[i].cnt = c
		i++
	}
	//sort.SortSlice(chars, func(i, j int) bool {return chars[i].cnt > chars[j].cnt})
	slices.SortFunc(chars, func(c1, c2 charBucket) int { return c2.cnt - c1.cnt })

	buf := make([]byte, len(s))
	idx := 0
	for i := range chars {
		b := bytes.Repeat([]byte{chars[i].elem}, chars[i].cnt)
		n := copy(buf[idx:], b)
		idx += n
	}

	return string(buf)
}
