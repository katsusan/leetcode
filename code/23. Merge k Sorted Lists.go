package code

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func mergeKLists(lists []*ListNode) *ListNode {
	if len(lists) == 0 {
		return nil
	}

	step := 1

	for step < len(lists) {
		for i := 0; i < len(lists)-step; i += step * 2 {
			lists[i] = mergeTwo(lists[i], lists[i+step])
		}
		step *= 2
	}

	return lists[0]
}

func mergeTwo(list1, list2 *ListNode) *ListNode {
	res := new(ListNode)
	pr := res

	for list1 != nil || list2 != nil {
		pr.Next = new(ListNode)
		pr = pr.Next

		if list1 == nil && list2 != nil {
			pr.Val = list2.Val
			list2 = list2.Next
			continue
		}

		if list2 == nil && list1 != nil {
			pr.Val = list1.Val
			list1 = list1.Next
			continue
		}

		if list1.Val < list2.Val {
			pr.Val = list1.Val
			list1 = list1.Next
		} else {
			pr.Val = list2.Val
			list2 = list2.Next
		}
	}
	pr = nil
	return res.Next
}

/* -> better solution for mergeTwoLists
func mergeTwo(l1 *ListNode,l2 *ListNode) *ListNode{
    point := &(ListNode{Val:0})
    head := point
    for l1 != nil && l2 != nil{
        if l1.Val <= l2.Val{
            point.Next = l1
            l1 = l1.Next
        } else{
            point.Next = l2
            l2 = l1
            l1 = point.Next.Next
        }
        point = point.Next
    }
    if l1 == nil{
        point.Next = l2
    } else{
        point.Next = l1
    }
    return head.Next
}*/
