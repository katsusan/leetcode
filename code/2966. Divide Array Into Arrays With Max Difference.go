package code

import "slices"

func divideArray(nums []int, k int) [][]int {
	slices.Sort(nums)
	arr := make([][]int, 0, len(nums)/3)
	for i := 0; i < len(nums)-2; i += 3 {
		if nums[i+2]-nums[i] > k {
			return [][]int{}
		}
		arr = append(arr, []int{nums[i], nums[i+1], nums[i+2]})
	}

	return arr
}
