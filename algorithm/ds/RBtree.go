package ds

//reference: https://zh.wikipedia.org/wiki/%E7%BA%A2%E9%BB%91%E6%A0%91
//红黑树特性properties:
//	1. Each node is either red or black.
//		每个节点非黑即白
//	2. The root is black.
//		约定根节点为黑
//	3. All leaves (NIL) are black.
//		所有叶子节点(NIL)为黑
//	4. If a node is red, then both its children are black.
//		红节点的子节点为黑
//	5. Every path from a given node to any of its descendant NIL nodes contains the same number of black nodes.
//		从任意节点到其每个叶子节点的路径所含黑节点数目相同 ⭐

///红黑树相对于AVL树来说，牺牲了部分平衡性以换取插入/删除操作时少量的旋转操作，整体来说性能要优于AVL树。

const (
	RED   = 0
	BLACK = 1
)

// RBTree describes a Red-black tree with root node
type RBTree struct {
	root *RBNode
}

// RBNode describes a Red-black tree node
type RBNode struct {
	lnode  *RBNode
	rnode  *RBNode
	parent *RBNode
	color  int8
	val    int32
}

func (node *RBNode) grandparent() *RBNode {
	if node.parent == nil {
		return nil
	}
	return node.parent.parent
}

func (node *RBNode) uncle() *RBNode {
	gp := node.grandparent()
	if gp == nil {
		return nil
	}

	if node.parent == gp.lnode {
		return gp.rnode
	}
	return gp.lnode
}

func (node *RBNode) sibling() *RBNode {
	if node.parent != nil {
		if node.parent.lnode == node {
			return node.rnode
		} else {
			return node.lnode
		}
	}

	return nil
}

/*
func rotateRight(node *RBNode) {
	gp := node.grandparent()
	par := node.parent
	rchild := node.rnode

	par.lnode = rchild

	if rchild != nil {

	}

}
*/
