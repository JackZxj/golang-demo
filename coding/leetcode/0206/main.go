package main

import (
	"fmt"
	"reflect"
)

func main() {
	tests := []struct {
		head   []int
		output []int
	}{
		{[]int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{[]int{1, 2}, []int{2, 1}},
		{[]int{}, []int{}},
	}
	for _, test := range tests {
		l1 := slice2List(test.head)
		got := reverseList(l1)
		l3 := list2Slice(got)
		if !reflect.DeepEqual(l3, test.output) {
			fmt.Printf("input %+v, got: %v\n", test, l3)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

//  Definition for singly-linked list.
type ListNode struct {
	Val  int
	Next *ListNode
}

func slice2List(s []int) *ListNode {
	header := &ListNode{}
	t := header
	for i := range s {
		t.Next = &ListNode{Val: s[i]}
		t = t.Next
	}
	return header.Next
}

func list2Slice(l *ListNode) []int {
	res := make([]int, 0)
	for l != nil {
		res = append(res, l.Val)
		l = l.Next
	}
	return res
}

// func reverseList(head *ListNode) *ListNode {
// 	if head == nil || head.Next == nil {
// 		return head
// 	}
// 	var secend *ListNode
// 	for head != nil {
// 		secend = &ListNode{Val: head.Val, Next: secend}
// 		head = head.Next
// 	}
// 	return secend
// }

func reverseList(head *ListNode) *ListNode {
	var pre *ListNode
	cur := head
	for cur != nil {
		tmp := cur.Next
		cur.Next = pre
		pre = cur
		cur = tmp
	}
	return pre
}
