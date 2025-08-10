package code

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func getDecimalValue(head *ListNode) int {
	listVal := 0
	for cur := head; cur != nil; cur = cur.Next {
		listVal = listVal*2 + cur.Val
	}
	return listVal
}
