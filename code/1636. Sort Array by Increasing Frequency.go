package code

import (
	"cmp"
	"slices"
)

func frequencySort1636(nums []int) []int {
	cntM := make(map[int]int)
	for i := range nums {
		cntM[nums[i]]++
	}

	slices.SortFunc(nums, func(a, b int) int {
		if cntM[a] != cntM[b] {
			return cmp.Compare(cntM[a], cntM[b])
		}

		return cmp.Compare(b, a)
	})

	return nums

	// type freqDesc struct {
	//     elem int
	//     freq int
	// }

	// freqSlc := make([]freqDesc, 0, 8)
	// for e, c := range cntM {
	//     freqSlc = append(freqSlc, freqDesc{
	//         elem: e,
	//         freq: c,
	//     })
	// }

	// slices.SortFunc(freqSlc, func(f1, f2 freqDesc) int {
	//     if f1.freq - f2.freq != 0 {
	//         return cmp.Compare(f1.freq, f2.freq)
	//     }

	//     return cmp.Compare(f2.elem, f1.elem)
	// })

	// r := make([]int, 0, len(nums))
	// for _, f := range freqSlc {
	//     r = append(r, repeat(f.elem, f.freq)...)
	// }

	// return r
}

// func repeat(elem, cnt int) []int {
//     s := make([]int, cnt)
//     for i := range s {
//         s[i] = elem
//     }
//     return s
// }
