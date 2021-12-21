package main

import (
	"fmt"
	"reflect"
)

func main() {
	tests := []struct {
		input  []int
		output []int
	}{
		{[]int{0, 1, 0, 3, 12}, []int{1, 3, 12, 0, 0}},
		{[]int{0}, []int{0}},
		{[]int{1}, []int{1}},
		{[]int{0, 0, 0, 1}, []int{1, 0, 0, 0}},
		{[]int{2, 2, 2, 1}, []int{2, 2, 2, 1}},
	}
	for i, test := range tests {
		moveZeroes(test.input)
		if !reflect.DeepEqual(test.input, test.output) {
			fmt.Printf("%d: %+v\n", i, test)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func moveZeroes(nums []int) {
	n := 0
	for i := range nums {
		if nums[i] != 0 {
			nums[n] = nums[i]
			n++
		}
	}
	for n < len(nums) {
		nums[n] = 0
		n++
	}
}
