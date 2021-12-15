package code


// use recursion
func findTargetSumWaysRecur(nums []int, target int) int {
    if len(nums) == 0 {
        if target == 0 {
            return 1
        }
        return 0
    }
    
    nlen := len(nums)
    case1 := findTargetSumWaysRecur(nums[:nlen-1], target-nums[nlen-1])  // last symbol is '+'
    case2 := findTargetSumWaysRecur(nums[:nlen-1], target+nums[nlen-1])  // last symbol is '-'
    return case1 + case2
}

