// binary search
package main

import "fmt"

func main() {
	tests := []struct {
		nums   []int
		target int
		output int
	}{
		{[]int{1, 3, 5, 7}, 0, -1},
		{[]int{1, 3, 5, 7}, 1, 0},
		{[]int{1, 3, 5, 7}, 2, -1},
		{[]int{1, 3, 5, 7}, 3, 1},
		{[]int{1, 3, 5, 7}, 5, 2},
		{[]int{1, 3, 5, 7}, 7, 3},
		{[]int{1, 3, 5, 7}, 9, -1},
		{[]int{1, 3, 5, 7, 9}, 0, -1},
		{[]int{1, 3, 5, 7, 9}, 4, -1},
		{[]int{1, 3, 5, 7, 9}, 5, 2},
		{[]int{1, 3, 5, 7, 9}, 7, 3},
		{[]int{1, 3, 5, 7, 9}, 9, 4},
		{[]int{1, 3, 5, 7, 9}, 11, -1},
	}
	for _, test := range tests {
		if got := search(test.nums, test.target); got != test.output {
			fmt.Printf("%+v, got: %d\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func search(nums []int, target int) int {
	// return bSearch(nums, target, 0, len(nums)-1)
	return forSearch(nums, target)
}

func forSearch(nums []int, target int) int {
	l := 0
	h := len(nums) - 1

	for l < h {
		t := (l + h) / 2
		if nums[t] == target {
			return t
		} else if nums[t] < target {
			l = t + 1
		} else {
			h = t - 1
		}
	}
	if nums[l] > target || nums[l] < target {
		return -1
	}
	return l
}

func bSearch(nums []int, target, l, h int) int {
	if target < nums[l] || target > nums[h] {
		return -1
	}
	if t := (h + l) / 2; nums[t] == target {
		return t
	} else if nums[t] < target {
		return bSearch(nums, target, t+1, h)
	} else {
		return bSearch(nums, target, l, t-1)
	}
}
