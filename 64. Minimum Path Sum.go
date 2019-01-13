package main

func minPathSum(grid [][]int) int {
	m := len(grid)
	if m == 0 {
		return 0
	}
	n := len(grid[0])
	if n == 0 {
		return 0
	}

	path := make([][]int, m)
	for i := range path {
		path[i] = make([]int, n)
	}

	path[0][0] = grid[0][0]

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if i == 0 && j > 0 {
				path[i][j] = path[i][j-1] + grid[i][j]
			}

			if j == 0 && i > 0 {
				path[i][j] = path[i-1][j] + grid[i][j]
			}

			if i > 0 && j > 0 {
				if path[i-1][j] < path[i][j-1] {
					path[i][j] = path[i-1][j] + grid[i][j]
				} else {
					path[i][j] = path[i][j-1] + grid[i][j]
				}
			}
		}
	}
	return path[m-1][n-1]
}
