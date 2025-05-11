package code

func updateMatrix(mat [][]int) [][]int {
	m, n := len(mat), len(mat[0])
	q := make([][2]int, 0, m*n/2)

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if mat[i][j] == 0 {
				q = append(q, [2]int{i, j})
			} else {
				mat[i][j] = m + n
			}
		}
	}

	for len(q) != 0 {
		cur := q[0]
		q = q[1:]
		adjs := adj(m, n, cur[0], cur[1])
		for _, adj := range adjs {
			if mat[adj[0]][adj[1]] > mat[cur[0]][cur[1]]+1 {
				mat[adj[0]][adj[1]] = mat[cur[0]][cur[1]] + 1
				q = append(q, [2]int{adj[0], adj[1]})
			}

		}
	}

	return mat
}

// for m x n matrix, return the adjacents of [i, j]
func adj(m, n, i, j int) [][2]int {
	a := make([][2]int, 0, 4)
	dmy := [][2]int{{i - 1, j}, {i + 1, j}, {i, j - 1}, {i, j + 1}}
	for _, v := range dmy {
		if v[0] >= 0 && v[0] < m && v[1] >= 0 && v[1] < n {
			a = append(a, v)
		}
	}
	return a
}
