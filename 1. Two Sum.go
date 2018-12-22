package main

//my solution ->70ms level
func twoSum(nums []int, target int) []int {
	var i, j int
	var result []int
	for i = 0; i < len(nums); i++ {
		for j = i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				result = []int{i, j}
			}
		}
	}
	return result
}

//better solution ->4ms level
/*
func twoSum(nums []int, target int) []int {
	m := make(map[int]int)

	for i := 0; i < len(nums); i++ {
		m[nums[i]] = i
	}
	ret := make([]int, 2)
	for i := 0; i < len(nums); i++ {

		a, x := m[target-nums[i]]
		if x && a != i {
			ret[0] = i
			ret[1] = a
			break
		}
	}
	return ret
}
*/
