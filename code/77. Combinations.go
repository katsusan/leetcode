package code

//use DP, for example: (5, 3) = (4, 2) + '5' && (4, 3)
//runtime: 244ms, top solution: 216ms
func combine(n int, k int) [][]int {
	var res [][]int
	if k <= 0 || n < k {
		return res
	}

	if n == k {
		var tmp []int
		for i := 1; i <= n; i++ {
			tmp = append(tmp, i)
		}
		res = append(res, tmp)
		return res
	}

	if k == 1 {
		for i := 1; i <= n; i++ {
			res = append(res, []int{i})
		}
		return res
	}

	for _, e := range combine(n-1, k-1) {
		res = append(res, append(e, n))
	}
	for _, e2 := range combine(n-1, k) {
		res = append(res, e2)
	}
	return res
}
