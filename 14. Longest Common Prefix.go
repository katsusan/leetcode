package main

//my solution (using recursion) ->20ms level
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	} else if len(strs) == 1 {
		return strs[0]
	} else if len(strs) == 2 {
		var res string
		var minlen int
		if len(strs[0]) < len(strs[1]) {
			minlen = len(strs[0])
		} else {
			minlen = len(strs[1])
		}

		for i := 0; i < minlen; i++ {
			if strs[0][i] == strs[1][i] {
				res += string(strs[0][i])
			} else {
				break
			}
		}
		return res
	}

	return longestCommonPrefix([]string{longestCommonPrefix(strs[0 : len(strs)-1]), strs[len(strs)-1]})

}

//better solution ->0ms level
/*
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	var current byte
	for i := 0; true; i++ {
		for j := 0; j < len(strs); j++ {
			if i >= len(strs[j]) {
				return strs[0][0:i]
			}

			if j == 0 {
				current = strs[j][i]
				continue
			}

			if strs[j][i] != current {
				return strs[0][0:i]
			}
		}
	}

	return ""
}
*/
