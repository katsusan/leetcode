package code

func construct2DArray(original []int, m int, n int) [][]int {
	if m*n != len(original) {
		return [][]int{}
	}

	r := make([][]int, m)
	for i := range r {
		r[i] = make([]int, n)
		for j := range r[i] {
			r[i][j] = original[i*n+j]
		}
	}
	return r
}
