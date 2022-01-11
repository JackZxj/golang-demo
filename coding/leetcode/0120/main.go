package main

import "fmt"

func main() {
	tests := []struct {
		triangle [][]int
		output   int
	}{
		{[][]int{{-10}}, -10},
		{[][]int{{2}, {3, 4}, {6, 5, 7}, {4, 1, 8, 3}}, 11},
		{[][]int{{2}, {3, 4}, {6, 5, 7}, {4, 1, 8, 3}, {-10, -11, -12, -13, -14}}, -1},
		{[][]int{{2}, {3, 4}, {6, -5, 7}, {4, 1, -118, 3}, {-10, -11, -12, -13, -14}}, -131},
		{[][]int{{1}, {1, 1}, {1, 1, 1}, {1, 1, 1, 1}, {1, 1, 1, 1, 1}}, 5},
		{[][]int{{0}, {0, 0}, {0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0, 0}}, 0},
		{[][]int{{-10000}, {-10000, -10000}, {-10000, -10000, -10000}, {-10000, -10000, -10000, -10000}, {-10000, -10000, -10000, -10000, -10000}}, -50000},
		{[][]int{{10000}, {10000, 10000}, {10000, 10000, 10000}, {10000, 10000, 10000, 10000}, {10000, 10000, 10000, 10000, 10000}}, 50000},
	}
	for i, test := range tests {
		got := minimumTotal(test.triangle)
		if got != test.output {
			fmt.Printf("%d# input: %+v, got: %v\n", i, test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

var res [][]int

func minimumTotal(triangle [][]int) int {
	res = [][]int{}
	output := 2000001
	for i := range triangle {
		tmp := make([]int, i+1)
		for j := range tmp {
			tmp[j] = output
		}
		res = append(res, tmp)
	}
	for i := range triangle[len(triangle)-1] {
		output = min(output, dp(triangle, len(triangle)-1, i))
	}
	return output
}

func dp(triangle [][]int, depth, index int) int {
	if 0 > index || index > depth {
		return 2000001
	}
	if depth == 0 {
		res[0][0] = triangle[0][0]
		return res[0][0]
	}
	if res[depth][index] == 2000001 {
		res[depth][index] = min(dp(triangle, depth-1, index-1), dp(triangle, depth-1, index)) + triangle[depth][index]
	}
	return res[depth][index]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
