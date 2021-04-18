package code

//my solution ->4ms level
func lengthOfLongestSubstring(s string) int {
	var arr [256]int = [256]int{-1}
	for k := range arr {
		arr[k] = -1
	}
	start, maxL := -1, 0
	for i := 0; i < len(s); i++ {
		if arr[s[i]] > start {
			start = arr[s[i]]
		}
		arr[s[i]] = i
		if i-start >= maxL {
			maxL = i - start
		}
	}
	return maxL

}
