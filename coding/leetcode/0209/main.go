package main

import (
	"fmt"
)

func main() {
	tests := []struct {
		nums   []int
		target int
		want   int
	}{
		{[]int{2, 3, 1, 2, 4, 3}, 7, 2},
		{[]int{1, 4, 4}, 4, 1},
		{[]int{1, 1, 1, 1, 1, 1, 1, 1}, 11, 0},
		{[]int{1, 1, 1, 1, 1, 1, 1, 1}, 8, 8},
		{[]int{100, 1, 1, 1, 1, 1, 1, 1, 1}, 8, 1},
		{[]int{6, 5, 4, 3, 2, 1}, 7, 2},
		{[]int{10, 3, 2}, 6, 1},
	}
	for _, test := range tests {
		got := minSubArrayLen(test.target, test.nums)
		if test.want != got {
			fmt.Printf("%+v, got: %d\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func minSubArrayLen(target int, nums []int) int {
	sum := 0
	result := len(nums) + 1
	left := 0
	check := false
	for right := 0; right < len(nums); right++ {
		if nums[right] >= target {
			return 1
		}
		sum += nums[right]
		for ; sum >= target && left <= right; left++ {
			sum -= nums[left]
			check = true
		}
		if check && result > right-left+2 {
			result = right - left + 2
			check = false
		}
		// fmt.Printf("%d-%d# sum: %d, res: %d, left: %d\n", right, nums[right], sum, result, left)
	}
	// fmt.Println("--------------------------")
	if result == len(nums)+1 {
		return 0
	}
	return result
}
