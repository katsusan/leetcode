package main

func searchInsert(nums []int, target int) int {
	//use binary search, O(n) = lg(n)

	if target > nums[len(nums)-1] {
		return len(nums)
	} else if target < nums[0] {
		return 0
	}

	for i, j := 0, len(nums); ; {

		if j == (i+1) && nums[i] < target && nums[j] > target {
			return j
		}

		if target < nums[(i+j)/2] {
			j = (i + j) / 2
			continue
		} else if target > nums[(i+j)/2] {
			i = (i + j) / 2
			continue
		} else {
			return (i + j) / 2
		}

	}
}
