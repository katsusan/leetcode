package code

/*
Given an array nums, return true if the array was originally sorted in non-decreasing order,
then rotated some number of positions (including zero). Otherwise, return false.

There may be duplicates in the original array.

Note: An array A rotated by x positions results in an array B of the same length such that
 B[i] == A[(i+x) % A.length] for every valid index i.
*/

func check(nums []int) bool {
	turnIdx := -1
	for i := 1; i < len(nums); i++ {
		if nums[i] < nums[i-1] {
			turnIdx = i
			break
		}
	}

	if turnIdx == -1 {
		return true
	}

	if nums[len(nums)-1] > nums[0] {
		return false
	}

	for j := turnIdx + 1; j < len(nums); j++ {
		if nums[j] < nums[j-1] {
			return false
		}
	}

	return true
}
