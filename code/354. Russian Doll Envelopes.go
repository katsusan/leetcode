package code

import (
	"sort"
)

// DP, n^2
func maxEnvelopesDP(envelopes [][]int) int {
	sort.Slice(envelopes, func(i, j int) bool {
		return envelopes[i][0] < envelopes[j][0] && envelopes[i][1] > envelopes[j][1]
	})

	dp := make([]int, len(envelopes))
	dp[0] = 1

	for i := 1; i < len(dp); i++ {
		for j := 0; j < i; j++ {
			if envelopes[i][1] > envelopes[j][1] {
				dp[i] = max354(dp[i], dp[j]+1)
			}
		}
	}

	return max354(dp...)
}

// patience sort, nlogn
func maxEnvelopesPatSort(envelopes [][]int) int {
	sort.Slice(envelopes, func(i, j int) bool {
		return envelopes[i][0] < envelopes[j][0] ||
			(envelopes[i][0] == envelopes[j][0] && envelopes[i][1] > envelopes[j][1])
	})

	piles := make([]int, len(envelopes))
	pileIdx := 0
	for i := 0; i < len(piles); i++ {
		poker := envelopes[i][1]

		left, right := 0, pileIdx
		for left < right {
			mid := (left + right) / 2
			if piles[mid] < poker {
				left = mid + 1
			} else {
				right = mid
			}
		}

		if left == pileIdx {
			pileIdx++
		}
		piles[left] = poker
	}

	return pileIdx
}

// length of arr is ensured to be >=1.
func max354(arr ...int) int {
	maxVal := arr[0]

	for i := 1; i < len(arr); i++ {
		if arr[i] > maxVal {
			maxVal = arr[i]
		}
	}

	return maxVal
}
