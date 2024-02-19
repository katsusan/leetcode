package code

func increasingTriplet(nums []int) bool {
	if len(nums) < 3 {
		return false
	}

	leftmin := make([]int, len(nums))  // holds [0, ... , i]'s min value
	rightmax := make([]int, len(nums)) // hp;ds [i, ..., len-1]'s max value

	leftmin[0] = nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] < leftmin[i-1] {
			leftmin[i] = nums[i]
		} else {
			leftmin[i] = leftmin[i-1]
		}
	}

	n := len(nums)
	rightmax[n-1] = nums[n-1]
	for j := n - 2; j >= 0; j-- {
		if nums[j] > rightmax[j+1] {
			rightmax[j] = nums[j]
		} else {
			rightmax[j] = rightmax[j+1]
		}
	}

	for k := 1; k < n-1; k++ {
		if leftmin[k] < nums[k] && nums[k] < rightmax[k] {
			return true
		}
	}
	return false
}
