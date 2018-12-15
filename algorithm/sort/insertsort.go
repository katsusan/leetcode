package sort

func InsertSort(nums []int) {
	for i := 1; i < len(nums); i++ {
		for j := i; j > 0; j-- { //ensure nums[i]'s left part are sorted,  just like nums[i] is inserted into nums[0,...,(i-1)]
			if nums[j-1] > nums[j] {
				nums[j-1], nums[j] = nums[j], nums[j-1]
			}
		}
	}
}
