package code

func combinationSum(candidates []int, target int) [][]int {
	var res [][]int

	mulindex := comb(candidates, target)
	for _, idx := range mulindex {
		if len(idx) == 0 {
			continue
		}
		var temp []int
		for i, elem := range idx {
			if elem == 0 {
				continue
			}
			for range make([]int, elem) {
				temp = append(temp, candidates[i])
			}
		}
		res = append(res, temp)
	}
	return res
}

func comb(cands []int, target int) [][]int {

	if len(cands) == 0 {
		return [][]int{}
	}
	if len(cands) == 1 {
		if target%cands[0] != 0 {
			return [][]int{{}}
		}
		return [][]int{{target / cands[0]}}
	}
	if target == 0 {
		return [][]int{make([]int, len(cands))}
	}

	var res [][]int
	for k := 0; k <= target/cands[0]; k++ {
		suf := comb(cands[1:], target-k*cands[0])
		for _, elem := range suf {
			if len(elem) != 0 {
				res = append(res, append([]int{k}, elem...))
			}

		}
	}
	return res
}
