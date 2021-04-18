package code

import "fmt"

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func removeNthFromEnd(head *ListNode, n int) *ListNode {

	if head == nil || n <= 0 {
		return head
	}

	var length int
	m := make(map[int]*ListNode)
	var tmpNode *ListNode = head
	for ; tmpNode != nil; length++ {
		m[length] = tmpNode
		tmpNode = tmpNode.Next
	}
	fmt.Println(length)
	if n > length {
		return head
	} else if n == 1 && length == 1 {
		return nil
	} else if n == 1 && length > 1 {
		m[length-2].Next = nil
		return head
	} else if n > length {
		return head
	} else if n == length {
		return m[1]
	} else if n < length {
		m[length-n-1].Next = m[length-n+1]
		return m[0]
	}

	return m[0]
}
