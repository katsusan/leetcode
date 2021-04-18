package code

//use a new stack to verify the match state of brackets
func isValid(s string) bool {
	length := len(s)

	if length == 0 {
		return true
	}

	if length%2 != 0 {
		return false
	}

	pairs := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}

	stack := make([]byte, length+1)
	bottom := 0
	top := 0      //top of brackets stack
	farthest := 0 //pointer to string s
	for {
		if s[farthest] == '(' || s[farthest] == '[' || s[farthest] == '{' {
			top++
			stack[top] = s[farthest]
			farthest++
		} else {
			if pairs[s[farthest]] == stack[top] {
				top = top - 1
				farthest++
			} else {
				return false
			}
		}

		if top == bottom {
			break
		}
		if farthest >= length && top != bottom {
			return false
		}
	}
	return isValid(s[farthest:])
}
