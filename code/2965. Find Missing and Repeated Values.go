package code

/*
You are given a 0-indexed 2D integer matrix grid of size n * n with values in the range [1, n2].
Each integer appears exactly once except a which appears twice and b which is missing. The task is to find the repeating and missing numbers a and b.

Return a 0-indexed integer array ans of size 2 where ans[0] equals to a and ans[1] equals to b.
*/

func findMissingAndRepeatedValues(grid [][]int) []int {
	// two primary formulas
	// a. 1 + 2 + 3 + ... + n = n(n+1)/2
	// b. 1^2 + 2^2 + 3^2 + ... + n^2 = n(n+1)(2n+1)/6
	sum, sqsum := 0, 0
	n := len(grid)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			cur := grid[i][j]
			sum += cur
			sqsum += cur * cur
		}
	}
	n = n * n
	diffSum := sum - (n*(n+1))/2
	diffSqsum := sqsum - (n*(n+1)*(2*n+1))/6
	a := (diffSum + diffSqsum/diffSum) / 2
	b := a - diffSum
	return []int{a, b}
}
