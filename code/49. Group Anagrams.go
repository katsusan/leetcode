package code

import "sort"

func groupAnagrams(strs []string) [][]string {
	anagramMap := make(map[string][]string)
	for _, s := range strs {
		var alphaOrd = alphaOrder(s)
		if v, ok := anagramMap[alphaOrd]; !ok {
			anagramMap[alphaOrd] = []string{s}
		} else {
			anagramMap[alphaOrd] = append(v, s)
		}
	}

	var res = make([][]string, 0, len(anagramMap))
	for _, ag := range anagramMap {
		res = append(res, ag)
	}

	return res
}

func alphaOrder(s string) string {
	b := []byte(s)
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	return string(b)
}
