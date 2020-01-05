package main

func rob(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	dp := make([]int, len(nums)) //dp[i] means the nums[0~i]'s maximum amount
	dp[0] = nums[0]

	for i := 1; i < len(nums); i++ {
		var step2, step3 int
		if i >= 2 {
			step2 = dp[i-2]
		}
		if i >= 3 {
			step3 = dp[i-3]
		}

		dp[i] = max(step2+nums[i], step3+nums[i-1])
	}
	return dp[len(nums)-1]

}

func max(num1, num2 int) int {
	if num1 > num2 {
		return num1
	}
	return num2
}
