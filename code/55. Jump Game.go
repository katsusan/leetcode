package code

//use predict algorithm, every time before you go,
//consider the step + nums[curr+step], it is maximum distance that you can reach.
func canJump(nums []int) bool {
	if len(nums) == 0 {
		return false
	}

	for i := 0; i < len(nums); {
		if nums[i] >= len(nums)-1-i {
			return true
		} else if nums[i] == 0 {
			return false
		}
		step, max := nums[i], 0
		for j := nums[i]; j > 0; j-- {
			if max < j+nums[i+j] {
				step = j
				max = j + nums[i+j]
			}
		}
		if max == 0 {
			return false
		} else {
			i += step
		}
	}
	return true
}
