package main

//use hash
func majorityElement(nums []int) int {
	cmap := make(map[int]int)
	length := len(nums)
	for i := 0; i <= length-1; i++ {
		_, found := cmap[nums[i]]
		if !found {
			cmap[nums[i]] = 1
			if length == 1 {
				return nums[i]
			}
			continue
		}
		//already exists
		cmap[nums[i]]++
		if cmap[nums[i]] > length/2 {
			return nums[i]
		}
	}
	return 0
}
