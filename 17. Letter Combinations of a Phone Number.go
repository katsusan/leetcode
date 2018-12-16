package main

func letterCombinations(digits string) []string {
	m := map[string]string{
		"2": "abc",
		"3": "def",
		"4": "ghi",
		"5": "jkl",
		"6": "mno",
		"7": "pqrs",
		"8": "tuv",
		"9": "wxyz",
	}

	var result []string

	if len(digits) == 0 {
		return result
	}

	if len(digits) == 1 {
		for _, ch := range m[digits] {
			result = append(result, string(ch))
		}
		return result
	}

	//use recursive
	//for example, result("234") =  result("23") x "4".
	//usually for loop
	rs := letterCombinations(digits[:len(digits)-1])
	for _, pre := range rs {
		for _, suf := range m[digits[len(digits)-1:]] {
			result = append(result, pre+string(suf))
		}
	}

	return result

}
