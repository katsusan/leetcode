package code

func minOperations(nums []int) int {
	ops := 0
	for i := 1; i < len(nums); i++ {
		if nums[i] <= nums[i-1] {
			new := nums[i-1] + 1
			ops += (new - nums[i])
			nums[i] = new
		}
	}

	return ops
}
