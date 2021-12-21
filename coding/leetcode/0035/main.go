package main

import "fmt"

func main() {
	tests := []struct {
		nums   []int
		target int
		expect int
	}{
		{[]int{1, 3, 5, 7}, 0, 0},
		{[]int{1, 3, 5, 7}, 1, 0},
		{[]int{1, 3, 5, 7}, 2, 1},
		{[]int{1, 3, 5, 7}, 3, 1},
		{[]int{1, 3, 5, 7}, 4, 2},
		{[]int{1, 3, 5, 7}, 5, 2},
		{[]int{1, 3, 5, 7}, 6, 3},
		{[]int{1, 3, 5, 7}, 7, 3},
		{[]int{1, 3, 5, 7}, 8, 4},
		{[]int{1, 3, 5}, 0, 0},
		{[]int{1, 3, 5}, 1, 0},
		{[]int{1, 3, 5}, 2, 1},
		{[]int{1, 3, 5}, 3, 1},
		{[]int{1, 3, 5}, 4, 2},
		{[]int{1, 3, 5}, 5, 2},
		{[]int{1, 3, 5}, 6, 3},
		{[]int{1}, 0, 0},
	}
	for _, test := range tests {
		if got := searchInsert(test.nums, test.target); got != test.expect {
			fmt.Printf("%+v, %d\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func searchInsert(nums []int, target int) int {
	// return forSearch(nums, target)
	return bSearch(nums, 0, len(nums)-1, target)
}

func forSearch(nums []int, target int) int {
	l := 0
	h := len(nums) - 1
	if target <= nums[l] {
		return 0
	}
	if target > nums[h] {
		return h + 1
	}
	for l < h {
		t := (l + h) / 2
		if nums[t] == target {
			return t
		} else if nums[t] > target {
			h = t
		} else {
			l = t + 1
		}
	}
	if nums[l] >= target {
		return l
	}
	return l + 1
}

func bSearch(nums []int, l, h, target int) int {
	if target < nums[l] {
		return l
	}
	if target > nums[h] {
		return h + 1
	}
	if t := (l + h) / 2; nums[t] == target {
		return t
	} else if nums[t] < target {
		return bSearch(nums, t+1, h, target)
	} else {
		return bSearch(nums, l, t-1, target)
	}
}
