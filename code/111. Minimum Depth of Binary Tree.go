package code

import "container/list"

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

func minDepthBFS(root *TreeNode) int {
	if root == nil {
		return 0
	}

	var depth = 1
	q := list.New()
	q.PushBack(root)

	for q.Len() != 0 {
		size := q.Len()
		for i := 0; i < size; i++ {
			cur := q.Front().Value.(*TreeNode)
			q.Remove(q.Front())
			if cur.Left == nil && cur.Right == nil {
				return depth
			}

			if cur.Left != nil {
				q.PushBack(cur.Left)
			}

			if cur.Right != nil {
				q.PushBack(cur.Right)
			}
		}

		depth++
	}

	return depth
}
