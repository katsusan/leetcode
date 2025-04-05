package code

func isArraySpecial1(nums []int) bool {
	for i := 0; i < len(nums)-1; i++ {
		if (nums[i]&1)^(nums[i+1]&1) == 0 {
			return false
		}
	}

	return true
}
