package main

import (
	"fmt"
	"reflect"
)

func main() {
	tests := []struct {
		l1     []int
		l2     []int
		output []int
	}{
		{[]int{1, 2, 4}, []int{1, 3, 4}, []int{1, 1, 2, 3, 4, 4}},
		{[]int{}, []int{}, []int{}},
		{[]int{}, []int{0}, []int{0}},
		{[]int{0}, []int{0}, []int{0, 0}},
		{[]int{0}, []int{1}, []int{0, 1}},
	}
	for _, test := range tests {
		l1 := slice2List(test.l1)
		l2 := slice2List(test.l2)
		got := mergeTwoLists(l1, l2)
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

// func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
// 	if list2 == nil {
// 		return list1
// 	}
// 	if list1 == nil {
// 		return list2
// 	}
// 	if list1.Val < list2.Val {
// 		list1.Next = mergeTwoLists(list1.Next, list2)
// 		return list1
// 	}
// 	list2.Next = mergeTwoLists(list1, list2.Next)
// 	return list2
// }

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	if list2 == nil {
		return list1
	}
	if list1 == nil {
		return list2
	}
	small, big := list1, list2
	if list1.Val > list2.Val {
		small, big = list2, list1
	}
	list1, list2 = small, small
	small = small.Next
	for small != nil && big != nil {
		if small.Val < big.Val {
			list2.Next = small
			list2 = small
			small = small.Next
		} else {
			list2.Next = big
			list2 = big
			big = big.Next
		}
	}
	if big != nil {
		list2.Next = big
	} else {
		list2.Next = small
	}
	return list1
}
