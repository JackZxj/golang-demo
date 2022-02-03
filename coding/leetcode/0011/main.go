package main

import (
	"fmt"
)

func main() {
	tests := []struct {
		height []int
		want   int
	}{
		{[]int{1, 8, 6, 2, 5, 4, 8, 3, 7}, 49},
		{[]int{1, 1}, 1},
		{[]int{4, 3, 2, 1, 4}, 16},
		{[]int{1, 2, 1}, 2},
		{[]int{1, 2, 2, 1}, 3},
		{[]int{1, 4, 4, 1}, 4},
		{[]int{1, 4, 4, 1, 1, 1}, 5},
	}
	for _, test := range tests {
		got := maxArea(test.height)
		if test.want != got {
			fmt.Printf("%+v, got: %d\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func maxArea(height []int) int {
	area := -1
	for left, right := 0, len(height)-1; left < right; {
		if height[left] < height[right] {
			area = max(area, height[left]*(right-left))
			for ; left < right; left++ {
				if height[left] < height[left+1] {
					left++
					break
				}
			}
		} else {
			area = max(area, height[right]*(right-left))
			for ; left < right; right-- {
				if height[right] < height[right-1] {
					right--
					break
				}
			}
		}
	}
	return area
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
