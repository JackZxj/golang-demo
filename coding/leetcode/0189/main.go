package main

import (
	"fmt"
	"reflect"
)

func main() {
	tests := []struct {
		nums   []int
		k      int
		expect []int
	}{
		{[]int{-1, -100, 3, 99}, 2, []int{3, 99, -1, -100}},
		{[]int{1, 2, 3, 4, 5, 6, 7}, 3, []int{5, 6, 7, 1, 2, 3, 4}},
		{[]int{-1, -100, 3, 99}, 6, []int{3, 99, -1, -100}},
		{[]int{-1}, 6, []int{-1}},
	}
	for i, test := range tests {
		rotate(test.nums, test.k)
		if !reflect.DeepEqual(test.expect, test.nums) {
			fmt.Printf("%d: %+v\n", i, test)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func rotate(nums []int, k int) {
	k = k % len(nums)
	t := append(nums[len(nums)-k:], nums[:len(nums)-k]...)
	copy(nums, t) // 神奇的copy
}

func rightRotate(nums *[]int, k int) {
	t := *nums
	k = k % len(t)
	t = append(t[len(t)-k:], t[:len(t)-k]...)
	*nums = t
}

func nrotate(nums []int, k int) {
	var length = len(nums)
	var out = make([]int, length)
	for i, v := range nums {
		out[(i+k)%length] = v
	}

	// nums[0] = out[0]
	for k, v := range out {
		nums[k] = v
	}
}
