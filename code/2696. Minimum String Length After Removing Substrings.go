package code

func minLength(s string) int {
	st := make(stack, 0, len(s))
	for i := 0; i < len(s); i++ {
		switch {
		case s[i] == 'B' && len(st) != 0 && st[len(st)-1] == 'A':
			st.Pop()
		case s[i] == 'D' && len(st) != 0 && st[len(st)-1] == 'C':
			st.Pop()
		default:
			st.Push(s[i])
		}
	}

	return len(st)
}

type stack []byte

func (s *stack) Push(elem byte) {
	*s = append(*s, elem)
}

func (s *stack) Pop() byte {
	n := len(*s)
	if n == 0 {
		panic("pop on empty stack")
	}
	top := (*s)[n-1]
	*s = (*s)[:n-1]
	return top
}
