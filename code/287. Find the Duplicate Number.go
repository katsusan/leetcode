package code

func findDuplicate(nums []int) int {
	// for [1, n] numbers, if there are duplicate ones,
	// then at least one cycle exists, Eg, [1,3,4,2,2] means 1->3->2->4->2.
	// Tortoise and Hare's meeting works here.
	fast, slow := nums[0], nums[0]
	for {
		slow = nums[slow]
		fast = nums[nums[fast]]

		if fast == slow {
			break
		}
	}

	slow = nums[0]
	for slow != fast {
		fast = nums[fast]
		slow = nums[slow]
	}
	return slow
}
