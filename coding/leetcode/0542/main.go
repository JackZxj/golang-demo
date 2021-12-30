package main

import (
	"fmt"
	"reflect"
)

func main() {
	tests := []struct {
		mat [][]int
		dst [][]int
	}{
		{[][]int{{0, 0, 0}, {0, 1, 0}, {0, 0, 0}}, [][]int{{0, 0, 0}, {0, 1, 0}, {0, 0, 0}}},
		{[][]int{{1, 1, 1}, {1, 0, 1}, {1, 1, 1}}, [][]int{{2, 1, 2}, {1, 0, 1}, {2, 1, 2}}},
		{[][]int{{0, 0, 0}, {0, 1, 0}, {1, 1, 1}}, [][]int{{0, 0, 0}, {0, 1, 0}, {1, 2, 1}}},
	}
	for i, test := range tests {
		got := updateMatrix(test.mat)
		if !reflect.DeepEqual(got, test.dst) {
			fmt.Printf("%d# input %+v, got: %v\n", i, test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func updateMatrix(mat [][]int) [][]int {
	bfs(mat, 0, 0)
	return mat
}

var (
	dx = [4]int{1, 0, 0, -1}
	dy = [4]int{0, 1, -1, 0}
)

func bfs(mat [][]int, m, n int) {
	if mat[m][n] != 0 {
		min := 10000
		for i := range dx {
			if 0 <= m+dx[i] && m+dx[i] < len(mat) && 0 <= n+dy[i] && n+dy[i] < len(mat[0]) {
				if mat[m+dx[i]][n+dy[i]] < min {
					min = mat[m+dx[i]][n+dy[i]]
				}
			}
		}
		mat[m][n] = min + 1
	}
}
