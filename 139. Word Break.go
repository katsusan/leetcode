package main

func wordBreak(s string, wordDict []string) bool {
	if len(s) == 0 {
		return true
	}

	wordmap := make(map[string]struct{}, len(wordDict))
	for _, word := range wordDict {
		wordmap[word] = struct{}{}
	}

	//dp[i] means the s[i-1] is reacheable or not
	dp := make([]bool, len(s)+1)
	dp[0] = true

	for i := 1; i <= len(s); i++ {
		for j := i - 1; j >= 0; j-- {
			_, ok := wordmap[s[j:i]]
			if dp[j] && ok {
				dp[i] = true
				continue
			}
		}
	}

	return dp[len(s)]
}
