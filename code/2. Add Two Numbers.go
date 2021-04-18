package code

//my solution  ->20ms level

//Definition for singly-linked list.
type ListNode struct {
	Val  int
	Next *ListNode
}

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {

	if l1 == nil {
		return l2
	} else if l2 == nil {
		return l1
	}

	var tmpl1, tmpl2, tmplr *ListNode
	resl := new(ListNode)
	tmplr = resl

	var ei, ss int = 0, 0

	for tmpl1, tmpl2 = l1, l2; tmpl1 != nil && tmpl2 != nil; tmpl1, tmpl2 = tmpl1.Next, tmpl2.Next {
		ss = tmpl1.Val + tmpl2.Val + ei
		if ss >= 10 {
			tmplr.Val = ss - 10
			ei = 1
		} else {
			tmplr.Val = ss
			ei = 0
		}
		if tmpl1.Next != nil && tmpl2.Next != nil {
			tmplr.Next = new(ListNode)
			tmplr = tmplr.Next
		}

	}
	if tmpl1 == nil && tmpl2 == nil && ei == 1 {
		tmplr.Next = new(ListNode)
		tmplr.Next.Val = ei
		return resl
	}

	if tmpl1 != nil {
		tmplr.Next = new(ListNode)
		tmplr = tmplr.Next
		for ; tmpl1 != nil; tmpl1 = tmpl1.Next {
			ss = tmpl1.Val + ei
			if ss == 10 {
				tmplr.Val = 0
				ei = 1
			} else {
				tmplr.Val = ss
				ei = 0
			}
			if tmpl1.Next != nil {
				tmplr.Next = new(ListNode)
				tmplr = tmplr.Next
			}
			if tmpl1.Next == nil && ei == 1 {
				tmplr.Next = new(ListNode)
				tmplr.Next.Val = ei
			}
		}
	}

	if tmpl2 != nil {
		tmplr.Next = new(ListNode)
		tmplr = tmplr.Next
		for ; tmpl2 != nil; tmpl2 = tmpl2.Next {
			ss = tmpl2.Val + ei
			if ss == 10 {
				tmplr.Val = 0
				ei = 1
			} else {
				tmplr.Val = ss
				ei = 0
			}
			if tmpl2.Next != nil {
				tmplr.Next = new(ListNode)
				tmplr = tmplr.Next
			}
			if tmpl2.Next == nil && ei == 1 {
				tmplr.Next = new(ListNode)
				tmplr.Next.Val = ei
			}
		}
	}
	return resl
}
