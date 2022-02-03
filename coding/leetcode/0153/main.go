package main

import "fmt"

func main() {
	tests := []struct {
		nums []int
		min  int
	}{
		{[]int{3}, 3},
		{[]int{3, 2}, 2},
		{[]int{3, 4, 2}, 2},
		{[]int{3, 4, 5, 1, 2}, 1},
		{[]int{2, 3, 4, 5, 1}, 1},
		{[]int{4, 5, 6, 7, 0, 1, 2}, 0},
		{[]int{11, 13, 15, 17}, 11},
	}
	for _, test := range tests {
		if got := findMin(test.nums); got != test.min {
			fmt.Printf("%+v, got: %v\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func findMin(nums []int) int {
	if nums[0] <= nums[len(nums)-1] {
		return nums[0]
	}
	return nums[findK(nums, 0, len(nums)-1)]
}

func findK(nums []int, start, end int) int {
	if start+1 == end {
		return end
	}
	mid := (start + end) / 2
	if nums[mid] < nums[end] {
		return findK(nums, start, mid)
	}
	return findK(nums, mid, end)
}
