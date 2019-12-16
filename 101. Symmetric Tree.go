package main

func isSymmetric(root *TreeNode) bool {
	if root == nil {
		return true
	}

	return isMirror(root.Left, root.Right)
}

func isMirror(left *TreeNode, right *TreeNode) bool {
	/*
		    if left != nil && right != nil {
		        return left.Val == right.Val && isMirror(left.Left, right.Right) && isMirror(left.Right, right.Left)
		    } else if left == nil && right == nil {
		        return true
			}
	*/

	if left == nil || right == nil {
		return left == right
	}

	return left.Val == right.Val && isMirror(left.Left, right.Right) && isMirror(left.Right, right.Left)
}
