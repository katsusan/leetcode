package main

//Description:
//Given an array nums and a value val, remove all instances of that value in-place and return the new length.
func removeElement(nums []int, val int) int {
	var cur int
	for cur = 0; cur < len(nums); cur++ {
		if nums[cur] == val {
			if idx := findEx(nums[cur+1:], val); idx == -1 {
				break
			} else {
				//exchange nums[cur] and nums[cur+idx+1]
				nums[cur], nums[cur+idx+1] = nums[cur+idx+1], nums[cur]
			}
		}
	}
	return cur
}

//return index of first element which is not equal to target.
//if all elements are equal to target, returns -1.
func findEx(nums []int, target int) int {
	for i := 0; i < len(nums); i++ {
		if nums[i] != target {
			return i
		}
	}
	return -1
}
