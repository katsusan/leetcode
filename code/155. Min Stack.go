package code

import "fmt"

type Node struct {
	Val  int
	Next *Node
}

type MinStack struct {
	TopNode *Node
	MinNode *Node
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

/** initialize your data structure here. */
func Constructor() MinStack {
	return MinStack{TopNode: nil, MinNode: nil}
}

func (this *MinStack) Push(x int) {
	if this.MinNode == nil {
		this.MinNode = &Node{Val: x, Next: nil}
	} else {
		fmt.Println("will extend ", min(x, this.MinNode.Val))
		this.MinNode = &Node{Val: min(x, this.MinNode.Val), Next: this.MinNode}
	}

	newTop := &Node{
		Val:  x,
		Next: this.TopNode,
	}
	this.TopNode = newTop

}

func (this *MinStack) Pop() {
	if this.TopNode != nil {
		this.TopNode = this.TopNode.Next
	}

	if this.MinNode != nil {
		this.MinNode = this.MinNode.Next
	}
}

func (this *MinStack) Top() int {
	if this.TopNode != nil {
		return this.TopNode.Val
	}
	return 0
}

func (this *MinStack) GetMin() int {
	var min int
	if this.MinNode != nil {
		min = this.MinNode.Val
	} else {
		min = 0
	}

	return min
}

/**
 * Your MinStack object will be instantiated and called as such:
 * obj := Constructor();
 * obj.Push(x);
 * obj.Pop();
 * param_3 := obj.Top();
 * param_4 := obj.GetMin();
 */
