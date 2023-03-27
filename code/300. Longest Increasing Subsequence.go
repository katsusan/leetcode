package code

// use DP, O(n^2) complexity
func lengthOfLIS(nums []int) int {
	dp := make([]int, len(nums))

	for i := 0; i < len(dp); i++ {
		dp[i] = 1
		for j := 0; j < i; j++ {
			if nums[i] > nums[j] {
				dp[i] = ArrMax(dp[i], dp[j]+1)
			}
		}
	}

	return ArrMax(dp...)
}

func ArrMax(arr ...int) int {
	maxVal := arr[0]
	for i := 1; i < len(arr); i++ {
		if arr[i] > maxVal {
			maxVal = arr[i]
		}
	}
	return maxVal
}

// patience sort, O(nlogn) complexity.
func lengthOFLIS2(nums []int) int {
	piles := make([]int, len(nums))
	pileCount := 0

	for i := 0; i < len(nums); i++ {
		pile := nums[i]
		left, right := 0, pileCount
		for left < right {
			mid := (left + right) / 2
			if pile <= piles[mid] {
				right = mid
			} else if pile > piles[mid] {
				left = mid + 1
			}
		}

		// currnem > all existing piles, get new pile
		if left == pileCount {
			pileCount++
		}

		piles[left] = pile
	}

	return pileCount
}
