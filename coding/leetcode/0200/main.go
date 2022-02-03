package main

import "fmt"

func main() {
	tests := []struct {
		grid [][]byte
		want int
	}{
		{[][]byte{{'1', '1', '1', '1', '0'}, {'1', '1', '0', '1', '0'}, {'1', '1', '0', '0', '0'}, {'0', '0', '0', '0', '0'}}, 1},
		{[][]byte{{'1', '1', '0', '0', '0'}, {'1', '1', '0', '0', '0'}, {'0', '0', '1', '0', '0'}, {'0', '0', '0', '1', '1'}}, 3},
	}
	for _, test := range tests {
		if got := numIslands(test.grid); got != test.want {
			fmt.Printf("%+v, got: %v\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

var (
	dx = [4]int{1, 0, 0, -1}
	dy = [4]int{0, 1, -1, 0}
)

func numIslands(grid [][]byte) int {
	count := 0
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == '1' {
				count++
				// fmt.Println("count", count, i, j)
				dfs(grid, i, j)
			}
		}
	}
	// fmt.Println("--------------------------")
	return count
}

func dfs(grid [][]byte, i, j int) {
	// fmt.Println(i, j)
	grid[i][j] = '0'
	for k := range dx {
		if 0 <= i+dx[k] && i+dx[k] < len(grid) &&
			0 <= j+dy[k] && j+dy[k] < len(grid[0]) &&
			grid[i+dx[k]][j+dy[k]] == '1' {
			dfs(grid, i+dx[k], j+dy[k])
		}
	}
}
