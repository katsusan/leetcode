package code

import "slices"

func largestLocal(grid [][]int) [][]int {
	size := len(grid) - 2
	m := make([][]int, size)
	for i := 0; i < size; i++ {
		m[i] = make([]int, size)
		for j := 0; j < size; j++ {
			m[i][j] = slices.Max(gridAround(grid, i+1, j+1))
		}
	}
	return m
}

func gridAround(grid [][]int, i, j int) []int {
	s := make([]int, 0, 9)
	for row := i - 1; row <= i+1; row++ {
		for col := j - 1; col <= j+1; col++ {
			s = append(s, grid[row][col])
		}
	}
	return s
}
