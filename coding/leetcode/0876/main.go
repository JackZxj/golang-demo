package main

import "fmt"

// Definition for singly-linked list.
type ListNode struct {
	Val  int
	Next *ListNode
}

func main() {
	tests := []struct {
		num    int // 节点数
		output int // 节点值，从1开始
	}{
		{1, 1},
		{5, 3},
		{5, 3},
		{6, 4},
		{100, 51},
		{120, 61},
	}
	for _, test := range tests {
		head := makeList(test.num)
		if got := middleNode(head); got.Val != test.output {
			fmt.Printf("%+v, got: %q\n", test, got.Val)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func makeList(num int) *ListNode {
	l := &ListNode{
		Val: 1,
	}
	head := l
	for i := 2; i <= num; i++ {
		ll := ListNode{
			Val: i,
		}
		l.Next = &ll
		l = &ll
	}
	return head
}

func middleNode(head *ListNode) *ListNode {
	count := 1
	mid := head
	for head.Next != nil {
		count++
		if count%2 == 0 {
			mid = mid.Next
		}
		head = head.Next
	}
	return mid
}
