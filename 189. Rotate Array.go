package main

//use temporary space to store rotated part.
func rotate(nums []int, k int) {
	if len(nums) <= 1 {
		return
	}
	realrt := k % len(nums)
	tmpslice := make([]int, realrt)
	copy(tmpslice, nums[len(nums)-realrt:])
	copy(nums[realrt:], nums[:len(nums)-realrt])
	copy(nums[:realrt], tmpslice)
}
