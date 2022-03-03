package code

func ConvertToTitle(n int) string {
	if n <= 26 {
		return string('A' + uint8(n) - 1)
	}
	rmd := n % 26
	if rmd == 0 {
		return ConvertToTitle((n-26)/26) + "Z"
	}
	return ConvertToTitle(n/26) + string('A'+uint8(rmd)-1)
}
