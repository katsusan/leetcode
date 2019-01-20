package main

//use bitmap to determine if every member of nums should be included.
//example:
//	nums = []int{1, 4, 6}
//	6 = binary(1 1 0) => means nums[0], nums[1] included and nums[2] excluded => []int{1, 4}
func subsets(nums []int) [][]int {
	length := len(nums)

	if length == 0 {
		return [][]int{}
	}

	var res [][]int

	for n := 0; n < (1 << uint16(length)); n++ {
		var sub []int
		for i := 0; i < length; i++ {
			index := 1 << uint(i)
			if n&index >= 1 {
				sub = append(sub, nums[i])
			}
		}
		res = append(res, sub)
	}
	return res
}
