package main

import (
	"fmt"
)

func main() {
	tests := []struct {
		nums []int
		k    int
		want int
	}{
		{[]int{10, 5, 2, 6}, 100, 8},
		{[]int{10, 5, 2, 6}, 1000, 10},
		{[]int{1, 2, 3}, 0, 0},
		{[]int{1, 2, 3}, 1, 0},
		{[]int{1, 2, 3}, 2, 1},
		{[]int{100, 99, 1, 2, 3}, 3, 3},
		{[]int{100, 99, 1, 2, 3}, 7, 6},
	}
	for _, test := range tests {
		got := numSubarrayProductLessThanK(test.nums, test.k)
		if test.want != got {
			fmt.Printf("%+v, got: %d\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func numSubarrayProductLessThanK(nums []int, k int) int {
	count := 0
	base := 1
	for start, end := 0, 0; start < len(nums); start++ {
		for ; end < len(nums); end++ {
			base *= nums[end]
			if base >= k {
				base /= nums[end]
				break
			}
		}
		if end-start > 0 {
			count += (end - start)
			base /= nums[start]
		} else {
			end = start + 1
		}
	}
	return count
}

func numSubarrayProductLessThanK2(nums []int, k int) int {
	count := 0
	for start := 0; start < len(nums); start++ {
		base := 1
		for end := start; end < len(nums); end++ {
			// fmt.Printf("%d: [%d,%d]\n", count, start, end)
			base *= nums[end]
			if base >= k {
				break
			}
			count++
		}
	}
	return count
}
