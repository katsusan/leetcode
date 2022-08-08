package code

import "sort"

// Given an integer array nums of size n, return the minimum number of
// moves required to make all array elements equal.

// In one move, you can increment or decrement an element of the array by 1.

func minMoves2(nums []int) int {
	sort.Ints(nums)

	mid := nums[len(nums)/2]
	moves := 0

	for _, v := range nums {
		if v < mid {
			moves += mid - v
		} else {
			moves += v - mid
		}
	}

	return moves
}
