package code

import "slices"

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func leafSimilar(root1 *TreeNode, root2 *TreeNode) bool {
	leaf1 := getLeaf(root1)
	leaf2 := getLeaf(root2)
	// fmt.Println("leaf1:", leaf1)
	// fmt.Println("leaf2:", leaf2)
	return slices.Equal(leaf1, leaf2)
}

func getLeaf(root *TreeNode) []int {
	if root == nil {
		return nil
	} else if root.Right == nil && root.Left == nil {
		return []int{root.Val}
	}

	return append(getLeaf(root.Left), getLeaf(root.Right)...)
}
