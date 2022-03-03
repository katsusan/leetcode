package ds

/*
*	reference: https://zh.wikipedia.org/wiki/%E7%BA%A2%E9%BB%91%E6%A0%91
*
*	红黑树特性properties:
*		1. Each node is either red or black.
*			每个节点非黑即白
*		2. The root is black.
*			约定根节点为黑
*		3. All leaves (NIL) are black.
*			所有叶子节点(NIL)为黑
*		4. If a node is red, then both its children are black.
*			红节点的子节点为黑
*		5. Every path from a given node to any of its descendant NIL nodes contains the same number of black nodes.
*			从任意节点到其每个叶子节点的路径所含黑节点数目相同 ⭐
*
*		红黑树相对于AVL树来说，牺牲了部分平衡性以换取插入/删除操作时少量的旋转操作，整体来说性能要优于AVL树。
*
*	插入case：	插入时我们令新节点初始为红色，这样可以尽量避免由于性质5导致的对树的调整，即插入红节点只会影响性质2和4
*
*		(1). 插入节点位于根上(即原树是空树)，此时直接令新节点为黑色即可
*		(2). 插入节点的父节点是黑，此时不违反性质2和4(？考虑新节点的子节点为红的情况)，无需操作
*		(3). 当前节点的父节点是红，且叔节点也为红，此时祖父节点一定存在。(此时不管父节点/当前节点是左还是右，对应方法都一样)
*			=>	将父节点和叔节点变为黑，祖父节点变为红，然后将当前节点指向祖父继续处理
*		(4). 当前节点的父节点是红色，叔节点为黑色，当前节点是右节点。
*			=>	以当前节点的父节点作为新的当前节点，并进行一次左旋
*		(5). 当前节点的父节点为红，叔节点为黑，当前节点为左节点。
*			=>	父节点变黑，祖节点变红，在组节点处右旋
*
*	删除case：	对于二叉查找树，删除一个节点时要么找出其左子树的最大值，要么找出右子树的最小值，将其移动到删除的节点处，
*				这样问题就化为删除至多一个儿子节点的问题。此时：
*				- 若删除的是一个红色节点，由于它的父节点和儿子都是黑色的，删除后可以直接用儿子节点来替换它
*				- 若删除的是黑色节点并且其儿子节点是红色，则删除后将其儿子顶替上来并涂为黑色
*				- 若删除的是黑色节点并且其儿子节点也是黑色，则分为以下几种：
*		(1). 如果新节点是根节点，则意味着我们从所有路径删除了一个黑节点，并且新根节点也为黑，不违反任何性质。/否则进入(2)
*		(2). 如果兄弟节点是红色，此时假设当前节点是父亲的左儿子(右儿子同理对调即可)，则以当前节点的父亲为中心进行一次左旋，
				然后对调父亲与祖父的颜色,这时候父节点是红色，兄弟节点是黑色，此时进入(4)/(5)/(6)
		(3). 如果父亲是黑色，兄弟和兄弟儿子也都是黑色，这时我们把兄弟节点涂为红色，这会导致通过当前父节点的所有路径都少一个黑色，
				因此要从父节点开始从(1)开始重新平衡。
			 否则进入(4)
		(4). 如果父亲是红色，但兄弟和兄弟儿子都是黑色，此时调换兄弟和父亲的颜色，可以保证补回了当前节点上删除的黑色节点，
			 否则进入(5)
		(5).
		(6).
*
*
*
*

*/

const (
	RED   = 0
	BLACK = 1
)

// RBTree describes a Red-Black tree with root node
type RBTree struct {
	root *RBNode
}

// RBNode describes a Red-Black tree node
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
