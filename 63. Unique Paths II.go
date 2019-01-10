package main

func uniquePathsWithObstacles(obstacleGrid [][]int) int {
	m := len(obstacleGrid)
	if m == 0 {
		return 0
	}
	n := len(obstacleGrid[0])
	if n == 0 {
		return 0
	}

	arr := make([][]int, m)
	for i := range arr {
		arr[i] = make([]int, n)
	}

	arr[0][0] = 1

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if obstacleGrid[i][j] == 1 {
				arr[i][j] = 0
				continue
			}
			if i == 0 && j > 0 {
				arr[i][j] = arr[i][j-1]
			}

			if j == 0 && i > 0 {
				arr[i][j] = arr[i-1][j]
			}

			if i > 0 && j > 0 {
				arr[i][j] = arr[i-1][j] + arr[i][j-1]
			}
		}
	}

	return arr[m-1][n-1]
}
