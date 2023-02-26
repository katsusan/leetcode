package code

func removeDuplicates(nums []int) int {
	var p1, p2 = 0, 0

	totalLen := len(nums)
	for p2 < totalLen {
		if nums[p2] > nums[p1] {
			nums[p1+1] = nums[p2]
			p1++
			continue
		}

		p2++
	}

	return p1 + 1
}
