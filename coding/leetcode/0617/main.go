package main

import (
	"fmt"
	"reflect"
)

// Definition for a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func main() {
	tests := []struct {
		root1  []int
		root2  []int
		expect []int
	}{
		{[]int{1, 3, 2, 5}, []int{2, 1, 3, 0, 4, 0, 7}, []int{3, 4, 5, 5, 4, 0, 7}},
		{[]int{1, 3, 2, 5}, []int{}, []int{1, 3, 2, 5}},
		{[]int{}, []int{2, 1, 3, 0, 4, 0, 7}, []int{2, 1, 3, 0, 4, 0, 7}},
	}
	for _, test := range tests {
		t1 := makeTree(test.root1)
		t2 := makeTree(test.root2)
		t3 := mergeTrees(t1, t2)
		l1 := makeList(t3)
		if !reflect.DeepEqual(l1, test.expect) {
			fmt.Printf("%+v, got: %v\n", test, l1)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func dfs1(list []int, index int) *TreeNode {
	if index <= len(list) && list[index-1] != 0 {
		return &TreeNode{Val: list[index-1], Left: dfs1(list, index<<1), Right: dfs1(list, (index<<1)+1)}
	}
	return nil
}

func makeTree(list []int) *TreeNode {
	return dfs1(list, 1)
}

func dfs2(root *TreeNode, list *[]int, index int) {
	if root != nil {
		if index > len(*list) {
			l := make([]int, index-len(*list))
			*list = append(*list, l...)
		}
		(*list)[index-1] = root.Val
		dfs2(root.Left, list, index<<1)
		dfs2(root.Right, list, (index<<1)+1)
	}
}

func makeList(root *TreeNode) []int {
	list := make([]int, 0)
	dfs2(root, &list, 1)
	return list
}

func mergeTrees(root1 *TreeNode, root2 *TreeNode) *TreeNode {
	if root1 == nil {
		return root2
	}
	if root2 == nil {
		return root1
	}
	root1.Val += root2.Val
	root1.Left = mergeTrees(root1.Left, root2.Left)
	root1.Right = mergeTrees(root1.Right, root2.Right)
	return root1
}
