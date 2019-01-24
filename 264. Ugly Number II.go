package main

//c2,c3,c5 means the max factors of 2/3/5.
func nthUglyNumber(n int) int {

	arr := make([]int, n)
	arr[0] = 1

	var c2, c3, c5 int

	for count := 1; count < n; count++ {
		arr[count] = minnums(arr[c2]*2, arr[c3]*3, arr[c5]*5)
		if arr[count] == arr[c2]*2 {
			c2++
		}
		if arr[count] == arr[c3]*3 {
			c3++
		}
		if arr[count] == arr[c5]*5 {
			c5++
		}
	}
	return arr[n-1]
}

func minnums(nums ...int) int {

	min := nums[0]

	for _, v := range nums {
		if v < min {
			min = v
		}
	}
	return min
}
