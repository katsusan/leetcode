package code

import "slices"

func sortPeople(names []string, heights []int) []string {
	imd := make([]int, len(names))
	for i := range imd {
		imd[i] = i
	}
	slices.SortFunc(imd, func(i, j int) int {
		return heights[j] - heights[i]
	})
	nwnames := make([]string, len(names))
	for i := range nwnames {
		nwnames[i] = names[imd[i]]
	}
	return nwnames
}
