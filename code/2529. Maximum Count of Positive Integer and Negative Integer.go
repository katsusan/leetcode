package code

func maximumCount(nums []int) int {
	negIdx, posIdx := -1, len(nums)

	for low, high := 0, len(nums)-1; low <= high; {
		mid := low + (high-low)/2
		if nums[mid] >= 0 {
			high = mid - 1
		} else {
			low = mid + 1
			negIdx = mid
		}
	}

	for low, high := 0, len(nums)-1; low <= high; {
		mid := low + (high-low)/2
		if nums[mid] <= 0 {
			low = mid + 1
		} else {
			high = mid - 1
			posIdx = mid
		}
	}
	return max(negIdx+1, len(nums)-posIdx)
}
