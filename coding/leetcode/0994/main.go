package main

import "fmt"

func main() {
	tests := []struct {
		grid   [][]int
		output int
	}{
		{[][]int{{2, 1, 1}, {1, 1, 0}, {0, 1, 1}}, 4},
		{[][]int{{2, 1, 1}, {0, 1, 1}, {1, 0, 1}}, -1},
		{[][]int{{0, 2}}, 0},
		{[][]int{{0, 1}}, -1},
	}
	for _, test := range tests {
		if got := orangesRotting(test.grid); got != test.output {
			fmt.Printf("%+v, got: %d\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

type Point struct {
	i, j, min int
}

var (
	di = [4]int{1, 0, 0, -1}
	dj = [4]int{0, 1, -1, 0}
)

func orangesRotting(grid [][]int) int {
	count1 := 0
	queue2 := make([]Point, 0)
	maxMinute := 0
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == 1 {
				count1++
			} else if grid[i][j] == 2 {
				queue2 = append(queue2, Point{i, j, 0})
			}
		}
	}
	if count1 == 0 {
		return 0
	}
	if len(queue2) == 0 {
		return -1
	}
	for len(queue2) != 0 {
		for k := range di {
			if 0 <= queue2[0].i+di[k] && queue2[0].i+di[k] < len(grid) &&
				0 <= queue2[0].j+dj[k] && queue2[0].j+dj[k] < len(grid[0]) &&
				grid[queue2[0].i+di[k]][queue2[0].j+dj[k]] == 1 {
				count1--
				grid[queue2[0].i+di[k]][queue2[0].j+dj[k]] = 2
				queue2 = append(queue2, Point{queue2[0].i + di[k], queue2[0].j + dj[k], queue2[0].min + 1})
				maxMinute = queue2[0].min + 1
			}
		}
		queue2 = queue2[1:]
	}
	if count1 != 0 {
		return -1
	}
	return maxMinute
}

func orangesRotting2(grid [][]int) int {
	slice1 := make([]Point, 0)
	queue2 := make([]Point, 0)
	maxMinute := 0
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == 1 {
				slice1 = append(slice1, Point{i, j, 0})
			} else if grid[i][j] == 2 {
				queue2 = append(queue2, Point{i, j, 0})
			}
		}
	}
	if len(slice1) == 0 {
		return 0
	}
	if len(queue2) == 0 {
		return -1
	}
	for len(queue2) != 0 {
		for k := range di {
			if 0 <= queue2[0].i+di[k] && queue2[0].i+di[k] < len(grid) &&
				0 <= queue2[0].j+dj[k] && queue2[0].j+dj[k] < len(grid[0]) &&
				grid[queue2[0].i+di[k]][queue2[0].j+dj[k]] == 1 {
				grid[queue2[0].i+di[k]][queue2[0].j+dj[k]] = 2
				queue2 = append(queue2, Point{queue2[0].i + di[k], queue2[0].j + dj[k], queue2[0].min + 1})
				maxMinute = queue2[0].min + 1
			}
		}
		queue2 = queue2[1:]
	}
	for i := range slice1 {
		if grid[slice1[i].i][slice1[i].j] == 1 {
			return -1
		}
	}
	return maxMinute
}
