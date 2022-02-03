package main

import (
	"fmt"
)

func main() {
	tests := []struct {
		nums   []int
		target int
		res    int
	}{
		{[]int{5}, 5, 0},
		{[]int{4}, 5, -1},
		{[]int{1, 3}, 0, -1},
		{[]int{1, 3}, 1, 0},
		{[]int{1, 3}, 2, -1},
		{[]int{1, 3}, 3, 1},
		{[]int{1, 3}, 4, -1},
		{[]int{5, 4}, 3, -1},
		{[]int{5, 4}, 4, 1},
		{[]int{5, 4}, 5, 0},
		{[]int{5, 4}, 6, -1},
		{[]int{4, 5, 6, 7, 0, 1, 2}, -1, -1},
		{[]int{4, 5, 6, 7, 0, 1, 2}, 0, 4},
		{[]int{4, 5, 6, 7, 0, 1, 2}, 1, 5},
		{[]int{4, 5, 6, 7, 0, 1, 2}, 2, 6},
		{[]int{4, 5, 6, 7, 0, 1, 2}, 3, -1},
		{[]int{4, 5, 6, 7, 0, 1, 2}, 4, 0},
		{[]int{4, 5, 6, 7, 0, 1, 2}, 5, 1},
		{[]int{4, 5, 6, 7, 0, 1, 2}, 6, 2},
		{[]int{4, 5, 6, 7, 0, 1, 2}, 7, 3},
		{[]int{4, 5, 6, 7, 0, 1, 2}, 8, -1},
		{[]int{5, 6, 7, 0, 1, 2, 4}, -1, -1},
		{[]int{5, 6, 7, 0, 1, 2, 4}, 0, 3},
		{[]int{5, 6, 7, 0, 1, 2, 4}, 1, 4},
		{[]int{5, 6, 7, 0, 1, 2, 4}, 2, 5},
		{[]int{5, 6, 7, 0, 1, 2, 4}, 3, -1},
		{[]int{5, 6, 7, 0, 1, 2, 4}, 4, 6},
		{[]int{5, 6, 7, 0, 1, 2, 4}, 5, 0},
		{[]int{5, 6, 7, 0, 1, 2, 4}, 6, 1},
		{[]int{5, 6, 7, 0, 1, 2, 4}, 7, 2},
		{[]int{5, 6, 7, 0, 1, 2, 4}, 8, -1},
	}
	for _, test := range tests {
		got := search(test.nums, test.target)
		if got != test.res {
			fmt.Printf("%+v,got: %d\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func search(nums []int, target int) int {
	if len(nums) == 1 {
		if target != nums[0] {
			return -1
		}
		return 0
	}
	// k == 0, 二分查找
	if nums[0] < nums[len(nums)-1] {
		return bsearch(nums, target, 0, len(nums)-1)
	}
	// k != 0, 以k分界二分查找
	k := bsearchK(nums, 0, len(nums)-1)
	// fmt.Println("k =", k)
	if nums[0] <= target {
		return bsearch(nums, target, 0, k)
	}
	return bsearch(nums, target, k+1, len(nums)-1)
}

func bsearchK(nums []int, low, high int) int {
	if nums[low] > nums[low+1] {
		return low
	}
	mid := (low + high) / 2
	if nums[mid] > nums[high] {
		return bsearchK(nums, mid, high)
	}
	return bsearchK(nums, low, mid)
}

func bsearch(nums []int, target, low, high int) int {
	if low > high {
		return -1
	}
	// fmt.Println("low:", low, "high:", high)
	mid := (low + high) / 2
	if nums[mid] < target {
		return bsearch(nums, target, mid+1, high)
	} else if nums[mid] > target {
		return bsearch(nums, target, low, mid-1)
	} else {
		return mid
	}
}
