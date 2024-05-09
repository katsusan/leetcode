package code

func asteroidCollision(asteroids []int) []int {
	st := make([]int, 0, 8)
	//st = append(st, asteroids[0])
	for i := 0; i < len(asteroids); i++ {
		for {
			if isStable(st, asteroids[i]) {
				st = append(st, asteroids[i])
				break
			}

			// either newcomer win or both explodes or newcomer explodes
			top := len(st) - 1
			if st[top] < 0 {
				st = append(st, asteroids[i])
				break
			}

			t := st[top] + asteroids[i]
			if t == 0 {
				// both explodes
				st = st[:top]
				break
			} else if t < 0 {
				// newcomer wins
				st = st[:top]
				continue
			} else {
				// newcomer loses and explodes
				break
			}
		}
	}
	return st
}

func isStable(arr []int, asize int) bool {
	if len(arr) == 0 {
		return true
	}

	n := len(arr)
	if arr[n-1]^asize < 0 {
		return false
	}

	return true
}
