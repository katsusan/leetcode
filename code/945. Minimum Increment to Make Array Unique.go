package code

import "slices"

func minIncrementForUnique(nums []int) int {
	slices.Sort(nums)

	moves := 0
	for i := 1; i < len(nums); i++ {
		if nums[i] <= nums[i-1] {
			new := nums[i-1] + 1
			moves += (new - nums[i])
			nums[i] = new
		}
	}

	return moves
}
