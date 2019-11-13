package main

func search(nums []int, target int) int {
	length := len(nums)
	if length == 0 {
		return -1
	}
	if nums[length-1] >= nums[0] {
		return binarySearch(nums, target)
	}

	//otherwise pivot exists
	pivot := searchForPivot(nums)
	if beforep := binarySearch(nums[:pivot+1], target); beforep != -1 {
		return beforep
	} else if afterp := binarySearch(nums[pivot+1:], target); afterp != -1 {
		return pivot + afterp + 1
	}
	return -1
}

//return pivot index for nums
func searchForPivot(nums []int) int {
	length := len(nums)
	if length == 2 {
		return 0
	} else if length == 3 {
		if nums[1] < nums[2] {
			return 0
		}
		return 1
	}

	//condition with length >= 3
	mid := length / 2
	if nums[mid] > nums[length-1] {
		//before pivot
		return mid + searchForPivot(nums[mid:])
	} else if nums[mid] < nums[length-1] {
		//after pivot
		return searchForPivot(nums[:mid+1])
	}
	return -1
}

func binarySearch(nums []int, target int) int {
	start, end := 0, len(nums)-1
	for start <= end {
		mid := (start + end) / 2
		if target < nums[mid] {
			end = mid - 1
			continue
		} else if target > nums[mid] {
			start = mid + 1
			continue
		}
		return mid
	}
	return end - start
}
