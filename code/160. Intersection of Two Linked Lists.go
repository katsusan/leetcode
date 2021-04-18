package code

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
//User HashTable, time: O(M+N)
//Maybe the minimum time complexity is O(M+N), as you can't decide
//if the current node is intersect point until you traverse the last node
//of the two single-linked lists.
func getIntersectionNode(headA, headB *ListNode) *ListNode {
	intermap := make(map[*ListNode]int)

	for pointA := headA; pointA != nil; pointA = pointA.Next {
		intermap[pointA] = 1
	}

	for pointB := headB; pointB != nil; pointB = pointB.Next {
		if _, found := intermap[pointB]; found {
			return pointB
		}
	}

	return nil
}
