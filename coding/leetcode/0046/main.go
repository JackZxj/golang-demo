package main

import (
	"fmt"
)

func main() {
	tests := []struct {
		input  []int
		output [][]int
	}{
		{[]int{1}, [][]int{{1}}},
		{[]int{0, 1}, [][]int{{0, 1}, {1, 0}}},
		{[]int{1, 2, 3}, [][]int{{1, 2, 3}, {1, 3, 2}, {2, 3, 1}, {2, 1, 3}, {3, 1, 2}, {3, 2, 1}}},
		{[]int{1, 2, 3, 4}, [][]int{{1, 2, 3, 4}, {1, 2, 4, 3}, {1, 3, 4, 2}, {1, 3, 2, 4}, {1, 4, 2, 3}, {1, 4, 3, 2}, {2, 3, 4, 1}, {2, 3, 1, 4}, {2, 4, 1, 3}, {2, 4, 3, 1}, {2, 1, 3, 4}, {2, 1, 4, 3}, {3, 4, 1, 2}, {3, 4, 2, 1}, {3, 1, 2, 4}, {3, 1, 4, 2}, {3, 2, 4, 1}, {3, 2, 1, 4}, {4, 1, 2, 3}, {4, 1, 3, 2}, {4, 2, 3, 1}, {4, 2, 1, 3}, {4, 3, 1, 2}, {4, 3, 2, 1}}},
	}
	for _, test := range tests {
		if got := permute(test.input); !isEqual(got, test.output) {
			fmt.Printf("input %+v, got: %v\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func isEqual(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	amap := make(map[string]bool)
	for i := range a {
		amap[slice(a[i])] = false
	}
	for i := range b {
		if ok, exist := amap[slice(b[i])]; exist {
			if ok {
				fmt.Println("already exist:", b[i])
				return false
			}
			amap[slice(b[i])] = true
		} else {
			fmt.Println("not exist:", b[i])
			return false
		}
	}
	return true
}

func slice(a []int) string {
	s := ""
	for i := range a {
		s = fmt.Sprintf("%s%d", s, a[i])
	}
	return s
}

var res [][]int

func permute(nums []int) [][]int {
	res = make([][]int, 0)
	forP(nums, []int{})
	return res
}

func forP(nums []int, now []int) {
	// fmt.Printf("nums: %v, now: %v\n", nums, now)
	if len(nums) == 0 {
		// fmt.Println("append", now)
		res = append(res, now)
		return
	}
	for i := len(nums); i > 0; i-- {
		forP(nums[1:], append(now, nums[0]))
		nums = append(nums[1:], nums[0])
	}
}
