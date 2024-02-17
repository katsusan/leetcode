package code

import "slices"

/*
Given a string s, reverse only all the vowels in the string and return it.

The vowels are 'a', 'e', 'i', 'o', and 'u', and they can appear in both lower and upper cases, more than once.

Input: s = "hello"
Output: "holle"
*/

func reverseVowels(s string) string {
	s1 := []byte(s)
	i, j := 0, len(s1)-1
	vowels := []byte{'A', 'E', 'I', 'O', 'U', 'a', 'e', 'i', 'o', 'u'}
	for i < j {
		if !slices.Contains(vowels, s1[i]) {
			i++
			continue
		}

		if !slices.Contains(vowels, s1[j]) {
			j--
			continue
		}

		s1[i], s1[j] = s1[j], s1[i]
		i++
		j--
	}
	return string(s1)
}
