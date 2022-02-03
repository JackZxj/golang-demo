package main

import "fmt"

func main() {
	tests := []struct {
		isConnected [][]int
		want        int
	}{
		{[][]int{{1, 1, 0},
			{1, 1, 0},
			{0, 0, 1}}, 2},
		{[][]int{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}, 3},
		{[][]int{{1, 0, 0, 1}, {0, 1, 1, 0}, {0, 1, 1, 1}, {1, 0, 1, 1}}, 1},
	}
	for i, test := range tests {
		if got := findCircleNum(test.isConnected); got != test.want {
			fmt.Printf("%d# %+v, got: %v\n", i, test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func findCircleNum(grid [][]int) int {
	count := 0
	for i := range grid {
		for j := i; j < len(grid[i]); j++ {
			if grid[i][j] == 1 {
				count++
				dfs(grid, i, j)
			}
		}
	}
	return count
}

func dfs(grid [][]int, i, j int) {
	grid[i][j] = 0
	grid[j][i] = 0
	for k := 0; k < len(grid); k++ {
		if grid[j][k] == 1 {
			dfs(grid, j, k)
		}
	}
}
