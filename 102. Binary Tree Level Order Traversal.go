package main

func levelOrder(root *TreeNode) [][]int {
	levelOrderNodes := nodesbylevel(root)
	var levelOrderVals [][]int

	for _, row := range levelOrderNodes {
		vals := make([]int, len(row))
		for i, node := range row {
			vals[i] = node.Val
		}
		levelOrderVals = append(levelOrderVals, vals)
	}
	return levelOrderVals
}

func nodesbylevel(root *TreeNode) [][]*TreeNode {
	res := make([][]*TreeNode, 0)

	if root == nil {
		return res
	}

	res = append(res, []*TreeNode{root})
	level := 0
	for {
		newlevel := make([]*TreeNode, 0)
		for _, node := range res[level] {
			if node.Left != nil {
				newlevel = append(newlevel, node.Left)
			}
			if node.Right != nil {
				newlevel = append(newlevel, node.Right)
			}
		}

		//leaf nodes
		if len(newlevel) == 0 {
			break
		}

		res = append(res, newlevel)
		level++
	}

	return res
}
