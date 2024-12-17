package code

func maxMoves(grid [][]int) int {
	rows := len(grid)
	cols := len(grid[0])
	//fmt.Printf("rows=%d, cols=%d\n", rows, cols)
	mins := make([][]int, rows)
	for i := range mins {
		mins[i] = make([]int, cols)
		mins[i][cols-1] = 0 // last row is definitely 0 moves
	}

	nextRows := func(rowIdx int) []int {
		switch rowIdx {
		case 0:
			return []int{0, 1}
		case rows - 1:
			return []int{rows - 1, rows - 2}
		default:
			return []int{rowIdx - 1, rowIdx, rowIdx + 1}
		}
	}

	for j := cols - 2; j >= 0; j-- {
		for i := 0; i < rows; i++ {
			nrs := nextRows(i)
			for _, nr := range nrs {
				if grid[i][j] < grid[nr][j+1] {
					mins[i][j] = max(mins[i][j], mins[nr][j+1]+1)
				}
			}
		}
	}

	maxmv := mins[0][0]
	for k := 1; k < rows; k++ {
		maxmv = max(maxmv, mins[k][0])
	}
	return maxmv
}
