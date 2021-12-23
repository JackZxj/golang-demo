package main

import (
	"fmt"
	"reflect"
)

// Definition for singly-linked list.
type ListNode struct {
	Val  int
	Next *ListNode
}

func main() {
	tests := []struct {
		head   []int
		n      int
		output []int
	}{
		{[]int{1, 2, 3, 4, 5}, 2, []int{1, 2, 3, 5}},
		{[]int{1}, 1, []int{}},
		{[]int{1, 2}, 1, []int{1}},
		{[]int{1, 2}, 2, []int{2}},
	}
	for _, test := range tests {
		res := removeNthFromEnd(sliceToList(test.head), test.n)
		got := listToSlice(res)
		if !reflect.DeepEqual(got, test.output) {
			fmt.Printf("%+v, got: %v\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func sliceToList(num []int) *ListNode {
	l := new(ListNode)
	t := l
	for _, n := range num {
		t.Next = &ListNode{
			Val: n,
		}
		t = t.Next
	}
	return l.Next
}

func listToSlice(l *ListNode) []int {
	res := make([]int, 0)
	for l != nil {
		res = append(res, l.Val)
		l = l.Next
	}
	return res
}

func removeNthFromEnd(head *ListNode, n int) *ListNode {
	h := &ListNode{Next: head}
	t := h
	head = h
	for head.Next != nil {
		n--
		if n < 0 {
			t = t.Next
		}
		head = head.Next
	}
	t.Next = t.Next.Next
	return h.Next
}
