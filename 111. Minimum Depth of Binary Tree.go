package main

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func minDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}

	if root.Left == nil && root.Right == nil {
		return 1
	}

	ml := minDepth(root.Left)
	mr := minDepth(root.Right)

	if root.Right == nil || (root.Left != nil && ml <= mr) {
		return ml + 1
	} else if root.Left == nil || (root.Right != nil && ml > mr) {
		return mr + 1
	}

	return mr + 1

}
