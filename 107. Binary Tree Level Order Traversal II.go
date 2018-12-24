package main

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func levelOrderBottom(root *TreeNode) [][]int {
	resv := make([][]int, 0) // reverse order of correct result
	res := make([][]int, 0)

	if root == nil {
		return res
	}

	tmpnodes := []*TreeNode{root}

	for {
		//traverse each level of nodes
		tmplist := make([]int, 0)
		for _, node := range tmpnodes {
			tmplist = append(tmplist, node.Val)
		}
		resv = append(resv, tmplist)

		clearflg := false

		for _, node := range tmpnodes {
			if !clearflg {
				tmpnodes = []*TreeNode{}
				clearflg = true
			}

			if node.Left != nil {
				tmpnodes = append(tmpnodes, node.Left)
			}
			if node.Right != nil {
				tmpnodes = append(tmpnodes, node.Right)
			}
		}

		if len(tmpnodes) == 0 {
			break
		}
	}

	for i := len(resv) - 1; i > -1; i-- {
		res = append(res, resv[i])
	}

	return res

}
