package code

import "slices"

func partitionArray1(nums []int, k int) int {
	minSubseq := 0
	slices.Sort(nums)
	for i := 0; i < len(nums); {
		j := i + 1
		for j < len(nums) {
			if nums[j]-nums[i] > k {
				//fmt.Printf("cnt++, i=%d, j=%d\n", i, j)
				minSubseq++
				break
			}
			j++
		}
		i = j
	}
	return minSubseq + 1
}

func partitionArray2(nums []int, k int) int {
	maxv, minv := slices.Max(nums), slices.Min(nums)
	bitmap := make([]bool, maxv-minv+1)
	partition := 1
	for i := range nums {
		bitmap[nums[i]-minv] = true
	}

	anchor := 0
	for i := 0; i < len(bitmap); i++ {
		if bitmap[i] && i-anchor > k {
			partition++
			anchor = i
		}
	}
	return partition
}
