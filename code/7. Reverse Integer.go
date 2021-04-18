package code

import "math"

//my solution ->4
func reverse(x int) int {
	var xa [20]int
	var rl int
	var flag bool = true
	var n int = 0

	if x < 0 {
		x = -x
		flag = false
	}
	for ; x > 0; n++ {
		xa[n] = x % 10
		x = x / 10
	}

	for i := 0; i < n; i++ {
		rl += xa[i] * int(math.Pow(10, float64(n-1-i)))
	}
	if flag == false && rl < 2147483648 {
		return -rl
	}
	if flag == true && rl < 2147483647 {
		return rl
	}
	return 0
}
