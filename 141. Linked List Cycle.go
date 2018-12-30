package main

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
//time: O(n)=>traverse all nodes , space: O(n) => extra map struct
func hasCycle(head *ListNode) bool {
	dupflag := false
	visit := make(map[*ListNode]int)
	tmpNode := head

	for tmpNode != nil {
		_, found := visit[tmpNode]
		if found {
			dupflag = true
			break
		} else {
			visit[tmpNode] = 1
		}
		tmpNode = tmpNode.Next
	}

	return dupflag
}

//use race with different speed.
//time: O(n), space: O(1)
func hasCycle2(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return false
	}

	faster := head.Next.Next
	slower := head

	for {
		if faster == nil || faster.Next == nil {
			return false
		}

		if faster == slower {
			return true
		}

		faster = faster.Next.Next
		slower = slower.Next
	}
}
