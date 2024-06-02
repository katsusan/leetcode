package code

func numSteps(s string) int {
	cnt := 0
	for len(s) != 1 && s[0] != 1 {
		n := len(s) - 1
		switch s[n] {
		case '0':
			s = s[:n]
		case '1':
			s = Adder(s)
		}
		cnt++
		//fmt.Printf("cnt=%d, s=%s\n", cnt, s)
	}
	return cnt
}

func Adder(s string) string {
	b := []byte(s)
	if s[len(b)-1] == '0' {
		b[len(b)-1] = '1'
	} else {
		for i := len(b) - 1; i >= 0; i-- {
			if b[i] == '0' {
				b[i] = '1'
				break
			} else {
				b[i] = '0'
				if i == 0 {
					b = append([]byte{'1'}, b...)
				}
			}
		}
	}

	return string(b)
}
