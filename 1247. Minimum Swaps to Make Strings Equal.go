package main

func minimumSwap(s1 string, s2 string) int {
	var countx, county int
	for i := range s1 {
		if s1[i] != s2[i] {
			if s1[i] == 'x' {
				countx++
			} else {
				county++
			}
		}
	}

	var swaps int
	if (countx-county)%2 != 0 {
		return -1
	}
	swaps = countx/2 + county/2 + 2*(countx%2)
	return swaps
}
