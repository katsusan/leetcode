package code

import "fmt"

//easy to understand solution from @mzchen
//store the previous max/min product.
func maxProduct(nums []int) int {
	max, min, res := nums[0], nums[0], nums[0]

	for i := 1; i < len(nums); i++ {
		if nums[i] < 0 {
			max, min = min, max
		}

		max = getmax(nums[i], nums[i]*max)
		min = getmin(nums[i], nums[i]*min)

		res = getmax(res, max)
		fmt.Println(res, max, min)
	}
	return res
}

func getmax(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func getmin(x, y int) int {
	if x < y {
		return x
	}
	return y
}
