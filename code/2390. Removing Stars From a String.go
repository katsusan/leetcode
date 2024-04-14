package code

func removeStars(s string) string {
	q := make([]byte, 0, 8)
	for i := range s {
		if s[i] != '*' {
			q = append(q, s[i])
		} else {
			q = q[:len(q)-1]
		}
	}

	return string(q)
}
