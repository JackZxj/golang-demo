package main

import (
	"fmt"
	"reflect"
)

func main() {
	tests := []struct {
		root   []int
		expect []int
	}{
		{[]int{1, 2, 3, 4, 5, 6, 7}, []int{1, -9999, 2, 3, -9999, 4, 5, 6, 7, -9999}},
		{[]int{1}, []int{1, -9999}},
		{[]int{}, []int{}},
	}
	for _, test := range tests {
		root := slice2Root(test.root)
		node := connect(root)
		got := expect2Slice(node)
		if !reflect.DeepEqual(got, test.expect) {
			fmt.Printf("%+v, got: %v\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

type Node struct {
	Val   int
	Left  *Node
	Right *Node
	Next  *Node
}

func dfs1(root []int, index int) *Node {
	if index <= len(root) {
		return &Node{
			Val:   root[index-1],
			Left:  dfs1(root, index<<1),
			Right: dfs1(root, (index<<1)+1),
		}
	}
	return nil
}

func slice2Root(root []int) *Node {
	return dfs1(root, 1)
}

func expect2Slice(root *Node) []int {
	res := make([]int, 0)
	if root == nil {
		return res
	}
	res = append(res, root.Val)
	for t := root.Next; t != nil; t = t.Next {
		res = append(res, t.Val)
	}
	res = append(res, -9999)
	left := expect2Slice(root.Left)
	res = append(res, left...)
	return res
}

func connect(root *Node) *Node {
	if root != nil {
		root.Next = nil
		bfs(root)
		return root
	}
	return nil
}

func bfs(parent *Node) {
	if parent.Left != nil {
		parent.Left.Next = parent.Right
		if parent.Next != nil {
			parent.Right.Next = parent.Next.Left
		}
		bfs(parent.Left)
		bfs(parent.Right)
	}
}
