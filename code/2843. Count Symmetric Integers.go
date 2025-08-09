package code

func countSymmetricIntegers(low int, high int) int {
	cnt := 0
	for i := low; i <= high; i++ {
		if i >= 10 && i <= 99 {
			if i/10 == i%10 {
				cnt++
			}
		} else if i >= 1000 && i <= 9999 {
			highsum := i/1000 + (i/100)%10
			lowsum := i%10 + (i/10)%10
			if highsum == lowsum {
				cnt++
			}
		}
	}
	return cnt
}
