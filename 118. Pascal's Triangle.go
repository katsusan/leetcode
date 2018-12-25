package main

import "fmt"

//pascal's triangle
func generate(numRows int) [][]int {

	if numRows == 0 {
		return [][]int{}
	}

	res := make([][]int, numRows)

	for r := 0; r < numRows; r++ {
		fmt.Println("row:", r)
		res[r] = make([]int, r+1)
		res[r][0], res[r][r] = 1, 1

		if r > 1 {
			for c := 1; c <= r/2; c++ {
				res[r][c], res[r][r-c] = res[r-1][c-1]+res[r-1][c], res[r-1][c-1]+res[r-1][c]
			}
		}
	}

	return res

}
