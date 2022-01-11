package main

import "fmt"

func main() {
	tests := []struct {
		input  []int
		output int
	}{
		{[]int{1}, 1},
		{[]int{1, 2}, 2},
		{[]int{1, 2, 3}, 4},
		{[]int{1, 6, 3}, 6},
		{[]int{1, 2, 3, 1}, 4},
		{[]int{1, 2, 3, 4}, 6},
		{[]int{4, 2, 3, 4}, 8},
		{[]int{4, 2, 3, 4, 1}, 8},
		{[]int{4, 2, 3, 4, 2}, 9},
		{[]int{1, 2, 3, 1, 6, 2, 3, 4, 5, 2, 3, 4}, 22},
		{[]int{0, 0, 0, 0, 0}, 0},
		{[]int{0, 0, 0, 1, 0, 0, 0, 0, 0}, 1},
		{[]int{0, 0, 0, 1, 1, 0, 0, 0, 0}, 1},
	}
	for i, test := range tests {
		if got := rob(test.input); got != test.output {
			fmt.Printf("%d# input: %+v, got: %v\n", i, test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

var res []int

func rob(nums []int) int {
	if len(nums) == 1 {
		return nums[0]
	}
	res = make([]int, len(nums))
	res[0] = nums[0]
	res[1] = nums[1]
	for i := 2; i < len(res); i++ {
		res[i] = -1
	}
	return max(dp(nums, len(nums)-1), dp(nums, len(nums)-2))
}

func dp(nums []int, n int) int {
	if n < 0 {
		return 0
	}
	if n > 1 && res[n] == -1 {
		res[n] = max(dp(nums, n-2), dp(nums, n-3)) + nums[n]
	}
	return res[n]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
