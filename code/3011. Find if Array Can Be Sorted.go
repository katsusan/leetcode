package code

func canSortArray(nums []int) bool {
	// separate nums into array clusters, every cluster
	// share the same number of set bits, clusters should satisfy:
	// cluster[i].Max <= cluster[i+1].Min

	cl := [4]int{} //clPrevMin, clPrevMax, clNextMin, clNextMax
	lastSetBit := 0
	for i := 0; i < len(nums); {
		bits := setbits(nums[i])
		if bits != lastSetBit {
			// new cluster, get current cluster's max and min
			cl[2], cl[3] = nums[i], nums[i]
			lastSetBit = bits
			for i < len(nums) && setbits(nums[i]) == lastSetBit {
				cl[2] = min(cl[2], nums[i])
				cl[3] = max(cl[3], nums[i])
				i++
			}
			if cl[1] > cl[2] {
				return false
			}
			cl[0], cl[1] = cl[2], cl[3]
		}
	}

	return true
}

func setbits(n int) int {
	cnt := 0
	for n > 0 {
		if n&1 != 0 {
			cnt++
		}
		n >>= 1
	}
	return cnt
}
