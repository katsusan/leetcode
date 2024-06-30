package code

// This solution can hanle cyclic cases
func findCenter(edges [][]int) int {
	cnt := make(map[int]int)
	for _, e := range edges {
		cnt[e[0]]++
		cnt[e[1]]++
	}

	n := len(cnt)
	for node := range cnt {
		if cnt[node] == n-1 {
			return node
		}
	}
	return -1
}
