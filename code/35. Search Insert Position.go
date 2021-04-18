package code

func searchInsert(nums []int, target int) int {
	//use binary search, O(n) = lg(n)

	i, j := 0, len(nums)-1
	for i <= j {

		if target < nums[(i+j)/2] {
			j = (i+j)/2 - 1
			continue
		} else if target > nums[(i+j)/2] {
			i = (i+j)/2 + 1
			continue
		} else {
			return (i + j) / 2
		}
	}
	return i
}
