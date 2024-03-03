package code

import "slices"

/*
You are given an integer array nums and an integer k.

In one operation, you can pick two numbers from the array whose sum equals k and remove them from the array.

Return the maximum number of operations you can perform on the array.
*/

func maxOperations(nums []int, k int) int {
	slices.Sort(nums)
	pairs := 0
	i, j := 0, len(nums)-1

	for i < j {
		s := nums[i] + nums[j]
		if s < k {
			i++
		} else if s > k {
			j--
		} else {
			pairs++
			i++
			j--
		}
	}
	return pairs
}
