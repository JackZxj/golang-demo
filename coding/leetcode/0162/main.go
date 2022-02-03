package main

import (
	"fmt"
)

func main() {
	tests := []struct {
		nums  []int
		peeks []int
	}{
		{[]int{1}, []int{0}},
		{[]int{1, 0}, []int{0}},
		{[]int{0, 1}, []int{1}},
		{[]int{1, 2, 3, 1}, []int{2}},
		{[]int{1, 2, 3, 4}, []int{3}},
		{[]int{4, 3, 2, 1}, []int{0}},
		{[]int{1, 2, 3, 4, 5}, []int{4}},
		{[]int{5, 4, 3, 2, 1}, []int{0}},
		{[]int{1, 2, 1, 3, 5, 6, 4}, []int{1, 5}},
		{[]int{1, 2, 7, 3, 5, 6}, []int{2, 5}},
		{[]int{21, 4, 74, 83, 6, 4, 7}, []int{0, 3, 6}},
	}
	for _, test := range tests {
		got := findPeakElement(test.nums)
		if !isPeek(test.peeks, got) {
			fmt.Printf("%+v,got: %d\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func isPeek(peeks []int, num int) bool {
	for i := range peeks {
		if peeks[i] == num {
			return true
		}
	}
	return false
}

func findPeakElement(nums []int) int {
	if len(nums) == 1 {
		return 0
	}
	low, high := 0, len(nums)-1
	for low+1 < high {
		mid := (low + high) / 2
		if nums[mid] < nums[mid+1] {
			low = mid
		} else {
			if nums[mid-1] < nums[mid] {
				return mid
			}
			high = mid
		}
	}
	if nums[low] < nums[high] {
		return high
	}
	return low
}
