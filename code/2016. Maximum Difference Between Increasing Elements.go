package code

func maximumDifference(nums []int) int {
	maxDiff := -1
	pfmin := nums[0]
	for i := 1; i < len(nums); i++ {
		pfmin = min(pfmin, nums[i-1])
		nd := nums[i] - pfmin
		if nd > 0 {
			maxDiff = max(maxDiff, nd)
		}
	}

	return maxDiff
}
