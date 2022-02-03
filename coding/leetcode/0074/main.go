package main

import "fmt"

func main() {
	tests := []struct {
		matrix [][]int
		target int
		exist  bool
	}{
		{[][]int{{1, 3, 5, 7}}, 0, false},
		{[][]int{{1, 3, 5, 7}}, 3, true},
		{[][]int{{1, 3, 5, 7}}, 4, false},
		{[][]int{{1, 3, 5, 7}}, 9, false},
		{[][]int{{1}, {3}}, 0, false},
		{[][]int{{1}, {3}}, 1, true},
		{[][]int{{1}, {3}}, 2, false},
		{[][]int{{1}, {3}}, 3, true},
		{[][]int{{1}, {3}}, 4, false},
		{[][]int{{1, 3}, {5, 7}, {9, 11}}, 2, false},
		{[][]int{{1, 3}, {5, 7}, {9, 11}}, 3, true},
		{[][]int{{1, 3}, {5, 7}, {9, 11}}, 4, false},
		{[][]int{{1, 3}, {5, 7}, {9, 11}}, 5, true},
		{[][]int{{1, 3}, {5, 7}, {9, 11}}, 6, false},
		{[][]int{{1, 3}, {5, 7}, {9, 11}}, 7, true},
		{[][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}}, 3, true},
		{[][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}}, 13, false},
	}
	for _, test := range tests {
		if searchMatrix(test.matrix, test.target) != test.exist {
			fmt.Printf("%+v, got: %v\n", test, !test.exist)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func searchMatrix(matrix [][]int, target int) bool {
	if target < matrix[0][0] {
		return false
	}
	line, mid := 0, 0
	start, end := 0, len(matrix)-1
	for start+1 < end {
		line = (start + end) / 2
		if matrix[line][0] > target {
			end = line - 1
		} else if matrix[line][0] < target {
			start = line
		} else {
			return true
		}
	}
	if target >= matrix[end][0] {
		line = end
	} else {
		line = start
	}
	start, end = 0, len(matrix[line])-1
	for start+1 < end {
		mid = (start + end) / 2
		if matrix[line][mid] > target {
			end = mid
		} else if matrix[line][mid] < target {
			start = mid
		} else {
			return true
		}
	}
	return matrix[line][start] == target || matrix[line][end] == target
}
