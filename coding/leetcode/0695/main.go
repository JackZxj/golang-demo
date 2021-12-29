package main

import "fmt"

func main() {
	tests := []struct {
		grid   [][]int
		output int
	}{
		{[][]int{{0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0}, {0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0}, {0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, {0, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0}, {0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 0}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0}, {0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0}, {0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0}}, 6},
		{[][]int{{0, 0, 0, 0, 0, 0, 0, 0}}, 0},
		{[][]int{{0, 0, 0, 1, 1, 0, 0, 0}}, 2},
		{[][]int{{0, 0, 0, 1, 1, 0, 0, 0}, {0, 0, 0, 1, 1, 0, 0, 0}}, 4},
	}
	for _, test := range tests {
		if got := maxAreaOfIsland(test.grid); got != test.output {
			fmt.Printf("%+v, got: %d\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

var (
	dx = [4]int{1, 0, 0, -1}
	dy = [4]int{0, 1, -1, 0}
)

func maxAreaOfIsland(grid [][]int) int {
	max := 0
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			s := dfs(grid, i, j)
			if s > max {
				max = s
			}
		}
	}
	return max
}

func dfs(grid [][]int, r, c int) int {
	if grid[r][c] == 1 {
		res := 1
		grid[r][c] = 0
		for i := range dx {
			if 0 <= r+dx[i] && r+dx[i] < len(grid) && 0 <= c+dy[i] && c+dy[i] < len(grid[0]) {
				res += dfs(grid, r+dx[i], c+dy[i])
			}
		}
		return res
	}
	return 0
}
