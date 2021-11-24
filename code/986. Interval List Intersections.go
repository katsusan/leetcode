package code

func intervalIntersection(firstList [][]int, secondList [][]int) [][]int {
	var res = make(twoDimenSlice, 0, (len(firstList)+len(secondList))/2)
	for i, j := 0, 0; i < len(firstList) && j < len(secondList); {
		e1 := max(firstList[i][0], secondList[j][0])
		e2 := min(firstList[i][1], secondList[j][1])
		res.Append([]int{e1, e2})
		if firstList[i][1] >= secondList[j][1] {
			j++
		} else {
			i++
		}
	}
	return [][]int(res)
}

type twoDimenSlice [][]int

func (s *twoDimenSlice) Append(x []int) {
	if x[0] > x[1] {
		return
	}
	*s = append(*s, []int{x[0], x[1]})
}

// func min(x, y int) int {
//     if x <= y {
//         return x
//     }
//     return y
// }

// func max(x, y int) int {
//     if x >= y {
//         return x
//     }
//     return y
// }
