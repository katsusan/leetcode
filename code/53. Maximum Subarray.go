package code

func maxSubArray(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	var maxsub int
	dp := make([]int, len(nums))
	maxsub, dp[0] = nums[0], nums[0]
	for i := 1; i < len(nums); i++ {
		if dp[i-1] < 0 {
			dp[i] = nums[i]
		} else {
			dp[i] = dp[i-1] + nums[i]
		}
		maxsub = max(maxsub, dp[i])
	}
	return maxsub
}
