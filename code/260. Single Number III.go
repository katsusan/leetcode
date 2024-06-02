package code

func singleNumber(nums []int) []int {
	xsum := 0
	for i := range nums {
		xsum ^= nums[i]
	}

	lb := xsum & (-xsum)

	x, y := 0, 0
	for i := range nums {
		if nums[i]&lb == 0 {
			x ^= nums[i]
		}
	}

	y = xsum ^ x
	return []int{x, y}
}
