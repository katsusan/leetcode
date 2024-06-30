package code

import "sort"

func getAncestors(n int, edges [][]int) [][]int {
	rEdges := make(map[int]set, n)
	for _, e := range edges {
		if rEdges[e[1]] == nil {
			rEdges[e[1]] = make(map[int]bool)
		}
		rEdges[e[1]].Add(e[0])
	}

	cache := make(map[int][]int)
	var bfs func(n int) []int
	bfs = func(n int) []int {
		anc := make(set)
		if rEdges[n] == nil {
			return []int{}
		}

		if cache[n] != nil {
			return cache[n]
		}

		for n := range rEdges[n] {
			anc.Add(n)
			for _, nn := range bfs(n) {
				anc.Add(nn)
			}
		}

		r := make([]int, 0, len(anc))
		for n := range anc {
			r = append(r, n)
		}
		sort.Ints(r)
		cache[n] = r
		return r
	}

	res := make([][]int, n)
	for i := range res {
		res[i] = bfs(i)
	}
	return res
}

type set map[int]bool

func (s set) Add(elem int) {
	s[elem] = true
}

func (s set) Exist(elem int) bool {
	return s[elem]
}
