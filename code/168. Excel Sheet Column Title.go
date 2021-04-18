package code

func convertToTitle(n int) string {
	if n <= 26 {
		return string('A' + n - 1)
	}
	rmd := n % 26
	if rmd == 0 {
		return convertToTitle((n-26)/26) + "Z"
	}
	return convertToTitle(n/26) + string('A'+rmd-1)
}
