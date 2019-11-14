package main

func searchRange(nums []int, target int) []int {
	startindex := getstartidx(nums, target)
	endindex := getendidx(nums, target)

	return []int{startindex, endindex}
}

func getstartidx(nums []int, target int) int {
	start, end := 0, len(nums)-1
	for start <= end {
		mid := (start + end) / 2
		if target <= nums[mid] {
			end = mid - 1
			continue
		} else if target > nums[mid] {
			start = mid + 1
			continue
		}
	}
	if start > len(nums)-1 || nums[start] != target {
		return -1
	}
	return start
}

func getendidx(nums []int, target int) int {
	start, end := 0, len(nums)-1
	for start <= end {
		mid := (start + end) / 2
		if target < nums[mid] {
			end = mid - 1
			continue
		} else if target >= nums[mid] {
			start = mid + 1
			continue
		}
	}
	if end < 0 || nums[end] != target {
		return -1
	}
	return end
}
