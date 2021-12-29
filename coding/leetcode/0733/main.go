package main

import (
	"fmt"
	"reflect"
)

func main() {
	tests := []struct {
		image    [][]int
		sr       int
		sc       int
		newColor int
		expect   [][]int
	}{
		{[][]int{{1, 1, 1}, {1, 1, 0}, {1, 0, 1}}, 1, 1, 2, [][]int{{2, 2, 2}, {2, 2, 0}, {2, 0, 1}}},
		{[][]int{{1, 1, 1, 1}, {1, 1, 0, 1}, {1, 0, 1, 0}}, 1, 1, 2, [][]int{{2, 2, 2, 2}, {2, 2, 0, 2}, {2, 0, 1, 0}}},
	}
	for _, test := range tests {
		if got := floodFill(test.image, test.sr, test.sc, test.newColor); !reflect.DeepEqual(got, test.expect) {
			fmt.Printf("%+v, got: %v\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func floodFill(image [][]int, sr int, sc int, newColor int) [][]int {
	if image[sr][sc] == newColor {
		return image
	}
	// return dfs1(image, sr, sc, newColor)
	// return bfs1(image, sr, sc, newColor)
	dfs3(image, sr, sc, newColor, image[sr][sc])
	return image
}

type Parent struct {
	r, c int
}

func dfs1(image [][]int, sr int, sc int, newColor int) [][]int {
	parents := make([]Parent, 0)
	parents = append(parents, Parent{sr, sc})
	target := image[sr][sc]
	for len(parents) > 0 {
		p := parents[len(parents)-1]
		image[p.r][p.c] = newColor
		// up
		if p.r+1 < len(image) && image[p.r+1][p.c] == target && image[p.r+1][p.c] != -1 {
			parents = append(parents, Parent{p.r + 1, p.c})
			continue
		}
		// down
		if p.r-1 >= 0 && image[p.r-1][p.c] == target && image[p.r-1][p.c] != -1 {
			parents = append(parents, Parent{p.r - 1, p.c})
			continue
		}
		// left
		if p.c-1 >= 0 && image[p.r][p.c-1] == target && image[p.r][p.c-1] != -1 {
			parents = append(parents, Parent{p.r, p.c - 1})
			continue
		}
		// right
		if p.c+1 < len(image[0]) && image[p.r][p.c+1] == target && image[p.r][p.c+1] != -1 {
			parents = append(parents, Parent{p.r, p.c + 1})
			continue
		}
		parents = parents[:len(parents)-1]
	}
	return image
}

func bfs1(image [][]int, sr int, sc int, newColor int) [][]int {
	parents := make([]Parent, 0)
	parents = append(parents, Parent{sr, sc})
	target := image[sr][sc]
	for len(parents) > 0 {
		p := parents[0]
		image[p.r][p.c] = newColor
		parents = parents[1:]
		// up
		if p.r+1 < len(image) && image[p.r+1][p.c] == target && image[p.r+1][p.c] != -1 {
			parents = append(parents, Parent{p.r + 1, p.c})
		}
		// down
		if p.r-1 >= 0 && image[p.r-1][p.c] == target && image[p.r-1][p.c] != -1 {
			parents = append(parents, Parent{p.r - 1, p.c})
		}
		// left
		if p.c-1 >= 0 && image[p.r][p.c-1] == target && image[p.r][p.c-1] != -1 {
			parents = append(parents, Parent{p.r, p.c - 1})
		}
		// right
		if p.c+1 < len(image[0]) && image[p.r][p.c+1] == target && image[p.r][p.c+1] != -1 {
			parents = append(parents, Parent{p.r, p.c + 1})
		}
	}
	return image
}

func dfs2(image [][]int, sr int, sc int, newColor int, target int) {
	if 0 <= sr && sr < len(image) && 0 <= sc && sc < len(image[0]) {
		if image[sr][sc] == target {
			image[sr][sc] = newColor
			dfs2(image, sr+1, sc, newColor, target) // up
			dfs2(image, sr-1, sc, newColor, target) // down
			dfs2(image, sr, sc-1, newColor, target) // left
			dfs2(image, sr, sc+1, newColor, target) // right
		}
	}
}

var op = []int{1, 0, -1, 0, 0, -1, 0, 1}

func dfs3(image [][]int, sr int, sc int, newColor int, target int) {
	if image[sr][sc] == target {
		image[sr][sc] = newColor
		for i := 0; i < 8; i += 2 {
			if 0 <= sr+op[i] && sr+op[i] < len(image) && 0 <= sc+op[i+1] && sc+op[i+1] < len(image[0]) {
				dfs3(image, sr+op[i], sc+op[i+1], newColor, target)
			}
		}
	}
}

// func bfs2(image [][]int, sr int, sc int, newColor int, target int) {
// 	if 0 <= sr && sr < len(image) && 0 <= sc && sc < len(image[0]) {
// 		if image[sr][sc] == target {
// 			image[sr][sc] = newColor
// 			dfs2(image, sr+1, sc, newColor, target) // up
// 			dfs2(image, sr-1, sc, newColor, target) // down
// 			dfs2(image, sr, sc-1, newColor, target) // left
// 			dfs2(image, sr, sc+1, newColor, target) // right
// 		}
// 	}
// }
