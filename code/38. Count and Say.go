package code

import (
	"bytes"
	"strconv"
)

func countAndSay(n int) string {
	fr := "1"
	for i := 1; i < n; i++ {
		fr = doRLE(fr)
	}
	return fr
}

func doRLE(s string) string {
	var b bytes.Buffer
	for i := 0; i < len(s); {
		cur := s[i]
		j := i
		for ; j < len(s); j++ {
			if s[j] != cur {
				break
			}
		}
		n := strconv.Itoa(j - i)
		b.WriteString(n)
		b.WriteByte(cur)
		i = j
	}

	return b.String()
}
