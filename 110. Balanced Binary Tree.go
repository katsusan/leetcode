package main

import "fmt"

func height(root *TreeNode) int {
	if root == nil {
		return 0
	}

	left := height(root.Left)
	right := height(root.Right)

	if right > left {
		return right + 1
	}

	return 1 + left
}

func isBalanced(root *TreeNode) bool {
	if root == nil {
		return true
	}

	left := height(root.Left)
	right := height(root.Right)
	fmt.Println(left, right)

	diff := left - right
	return diff <= 1 && diff >= -1 && isBalanced(root.Left) && isBalanced(root.Right)
}
