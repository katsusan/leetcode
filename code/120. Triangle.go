package code

func minimumTotal(triangle [][]int) int {
	m := len(triangle)
	if m == 0 {
		return 0
	}
	for row := range triangle {
		if len(triangle[row]) != row+1 {
			return 0
		}
	}

	arr := make([]int, m)
	mir := make([]int, m)
	arr[0] = triangle[0][0]

	for i := 1; i < m; i++ {
		copy(mir, arr)
		arr[0] += triangle[i][0]
		arr[i] = mir[i-1] + triangle[i][i]
		for j := 1; j < i; j++ {
			if mir[j-1] < mir[j] {
				arr[j] = mir[j-1] + triangle[i][j]
			} else {
				arr[j] = mir[j] + triangle[i][j]
			}
		}
	}

	min := arr[0]
	for _, v := range arr {
		if v < min {
			min = v
		}
	}
	return min
}
