package main

import (
	"fmt"
	"reflect"
)

func main() {
	tests := []struct {
		nums   []int
		target int
		res    []int
	}{
		{[]int{}, 0, []int{-1, -1}},
		{[]int{5}, 5, []int{0, 0}},
		{[]int{4}, 5, []int{-1, -1}},
		{[]int{5, 7, 7, 8, 9, 10}, 8, []int{3, 3}},
		{[]int{5, 6, 7, 8, 9, 10}, 7, []int{2, 2}},
		{[]int{5, 7, 7, 8, 8, 10}, 8, []int{3, 4}},
		{[]int{5, 7, 7, 8, 8, 10}, 6, []int{-1, -1}},
		{[]int{5, 7, 8, 8, 8, 10}, 8, []int{2, 4}},
		{[]int{0, 0, 0, 0, 0, 0, 0, 0, 0}, 0, []int{0, 8}},
		{[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 0, []int{0, 9}},
		{[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, 0, []int{0, 9}},
		{[]int{-1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 0, []int{1, 10}},
	}
	for _, test := range tests {
		got := searchRange(test.nums, test.target)
		if !reflect.DeepEqual(got, test.res) {
			fmt.Printf("%+v,got: %d\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func searchRange(nums []int, target int) []int {
	res := []int{-1, -1}
	if len(nums) == 0 {
		return res
	}
	res[0] = getLow(nums, target, 0, len(nums)-1)
	if res[0] != -1 {
		res[1] = getHigh(nums, target, res[0], len(nums)-1)
	}
	return res
}

func getLow(nums []int, target, start, end int) int {
	if start+1 >= end {
		if nums[start] == target {
			return start
		} else if nums[end] == target {
			return end
		}
		return -1
	}
	mid := (start + end) / 2
	if nums[mid] >= target {
		return getLow(nums, target, start, mid)
	}
	return getLow(nums, target, mid, end)
}

func getHigh(nums []int, target, start, end int) int {
	if start+1 >= end {
		if nums[end] == target {
			return end
		} else if nums[start] == target {
			return start
		}
		return -1
	}
	mid := (start + end) / 2
	if nums[mid] <= target {
		return getHigh(nums, target, mid, end)
	}
	return getHigh(nums, target, start, mid)
}
