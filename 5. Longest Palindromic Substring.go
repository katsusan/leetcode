package main

//expand string from Palindrome's center
//n chars with each char using 2n times of executing. so totally time compelxity is n*n.
func longestPalindrome(s string) string {
	if len(s) <= 1 {
		return s
	}

	maxs := ""
	maxlen := 0

	for n := 0; n < len(s)-1; n++ {
		s_odd := expandCenter(s, n, n)
		s_even := expandCenter(s, n, n+1)
		if maxlen < len(s_odd) {
			maxs = s_odd
			maxlen = len(s_odd)
		}
		if maxlen < len(s_even) {
			maxs = s_even
			maxlen = len(s_even)
		}
	}

	return maxs
}

func expandCenter(s string, i, j int) string {
	for i >= 0 && j < len(s) && s[i] == s[j] {
		i--
		j++
	}
	i++
	j--
	return s[i : j+1]
}
