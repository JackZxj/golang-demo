package main

import (
	"fmt"
	"reflect"
)

func main() {
	tests := []struct {
		nums   []int
		expect []int
	}{
		{[]int{-4, -1, 0, 3, 10}, []int{0, 1, 9, 16, 100}},
		{[]int{-7, -3, 2, 3, 11}, []int{4, 9, 9, 49, 121}},
		{[]int{0}, []int{0}},
		{[]int{2, 2, 2, 2}, []int{4, 4, 4, 4}},
		{[]int{1, 2, 3, 4}, []int{1, 4, 9, 16}},
		{[]int{-4, -3, -2, -1}, []int{1, 4, 9, 16}},
	}
	for _, test := range tests {
		if got := sortedSquares(test.nums); !reflect.DeepEqual(got, test.expect) {
			fmt.Printf("%+v, %v", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func sortedSquares(nums []int) []int {
	l := 0
	h := len(nums) - 1
	res := make([]int, h+1)
	i := h
	for ; l <= h; i-- {
		if nums[h] * nums[h] > nums[l] * nums[l] {
			res[i] = nums[h] * nums[h]
			h--
		} else {
			res[i] = nums[l] * nums[l]
			l++
		}
	}
	return res
}

func sortedSquares0(nums []int) []int {
	l := 0
	h := len(nums) - 1
	res := make([]int, h+1)
	ll := nums[l] * nums[l]
	hh := nums[h] * nums[h]
	i := h
	for ; l < h; i-- {
		if hh > ll {
			res[i] = hh
			h--
			hh = nums[h] * nums[h]
		} else {
			res[i] = ll
			l++
			ll = nums[l] * nums[l]
		}
	}
	res[i] = ll
	return res
}