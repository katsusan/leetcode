package main

//my solution
func twoSum(nums []int, target int) []int {
	m := make(map[int]int)

	for i := range nums {
		m[nums[i]] = i
	}

	for j := range nums {
		if v, ok := m[target-nums[j]]; ok && j != v {
			return []int{j, v}
		}
	}
	return []int{}
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
