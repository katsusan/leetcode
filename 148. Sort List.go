package main

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func sortList(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	firstHalf, secondHalf := splitIntoHalves(head)
	return mergeList(sortList(firstHalf), sortList(secondHalf))
}

func splitIntoHalves(head *ListNode) (*ListNode, *ListNode) {
	fast, slow := head, head
	var slowprev *ListNode
	for fast != nil && fast.Next != nil {
		slowprev = slow
		slow = slow.Next
		fast = fast.Next.Next
	}
	slowprev.Next = nil

	return head, slow
}

func mergeList(listA *ListNode, listB *ListNode) *ListNode {
	res := &ListNode{}

	cur := res
	for listA != nil && listB != nil {
		if listA.Val < listB.Val {
			cur.Next = listA
			listA = listA.Next
		} else {
			cur.Next = listB
			listB = listB.Next
		}
		cur = cur.Next
	}

	if listA != nil {
		cur.Next = listA
	} else if listB != nil {
		cur.Next = listB
	}
	return res.Next
}
