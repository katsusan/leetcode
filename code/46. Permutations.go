package code

func permute(nums []int) [][]int {
    res := make([][]int, 0, factorial(len(nums)))
    backtrack(nums, []int{}, &res)
    return res
}

func backtrack(remains []int, track []int, res *[][]int) {
    if len(remains) == 0 {
        *res = append(*res, track)
        return
    }
    
    for i := range remains {
        backtrack(subSlice(remains, i), append(track, remains[i]), res)
    }
}

// subSLice remove the element of s[index] and doesn't change original slice.
func subSlice(s []int, index int) []int {
    if index < 0 || index >= len(s) || len(s) == 0 {
        return []int{}
    }
    
    res := make([]int, len(s)-1)
    copy(res, s[:index])
    copy(res[index:], s[index+1:])
    
    return res
}

func factorial(n int) int {
    var fact int
    for i := 1; i <= n; i++ {
        fact *= i
    }
    return fact
}