package main

import (
	"fmt"
	"reflect"
)

func main() {
	tests := []struct {
		head []int
		want []int
	}{
		{[]int{1}, []int{1}},
		{[]int{1, 1}, nil},
		{[]int{1, 1, 3}, []int{3}},
		{[]int{1, 1, 3, 3}, nil},
		{[]int{1, 1, 3, 3, 4}, []int{4}},
		{[]int{1, 1, 1, 1, 4}, []int{4}},
	}
	for _, test := range tests {
		t1 := slice2list(test.head)
		t2 := deleteDuplicates(t1)
		got := list2slice(t2)
		if !reflect.DeepEqual(got, test.want) {
			fmt.Printf("%+v,got: %v\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

// Definition for singly-linked list.
type ListNode struct {
	Val  int
	Next *ListNode
}

func list2slice(list *ListNode) []int {
	var res []int
	for list != nil {
		res = append(res, list.Val)
		list = list.Next
	}
	return res
}

func slice2list(slice []int) *ListNode {
	list := &ListNode{}
	head := list
	for _, n := range slice {
		list.Next = &ListNode{Val: n}
		list = list.Next
	}
	return head.Next
}

func deleteDuplicates(head *ListNode) *ListNode {
	head = &ListNode{Next: head}
	pre, now := head, head.Next
	for now != nil && now.Next != nil {
		if now.Val == now.Next.Val {
			for now.Next != nil && now.Val == now.Next.Val {
				now.Next = now.Next.Next
			}
			pre.Next = now.Next
		} else {
			pre = now
		}
		now = now.Next
	}
	return head.Next
}
