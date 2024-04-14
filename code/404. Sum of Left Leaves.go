package code

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func sumOfLeftLeaves(root *TreeNode) int {
	if root == nil {
		return 0
	}

	s := 0
	traverse404(root.Left, root, &s)
	traverse404(root.Right, root, &s)
	return s
}

func traverse404(node, parent *TreeNode, sum *int) {
	// left leaf must have a parent
	if node == nil {
		return
	}

	// leaf must haven't any children
	if node.Left == nil && node.Right == nil && node == parent.Left {
		*sum += node.Val
	}

	traverse404(node.Left, node, sum)
	traverse404(node.Right, node, sum)
}
