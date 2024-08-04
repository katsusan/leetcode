package code

import "slices"

func rangeSum(nums []int, n int, left int, right int) int {
	sums := make([]int, 0, n*(n+1)/2)

	for i := 0; i < n; i++ {
		s := 0
		for j := i; j < n; j++ {
			s += nums[j]
			sums = append(sums, s)
		}
	}

	slices.Sort(sums)
	r := 0
	m := int(1e9 + 7)
	//fmt.Println("sums:", sums)
	for i := left - 1; i < right; i++ {
		if sums[i] >= m {
			r += sums[i] % m
		} else {
			r += sums[i]
		}

		if r >= m {
			r = r % m
		}
	}

	return r
}
