package main

import (
	"fmt"
	"reflect"
)

func main() {
	tests := []struct {
		nums   []int
		target int
		output []int
	}{
		{[]int{2, 7, 11, 15}, 9, []int{1, 2}},
		{[]int{2, 2, 11, 15}, 4, []int{1, 2}},
		{[]int{2, 3, 4}, 6, []int{1, 3}},
		{[]int{-1, 0}, -1, []int{1, 2}},
		{[]int{0, 3, 6, 7}, 10, []int{2, 4}},
	}
	for _, test := range tests {
		if got := twoSum(test.nums, test.target); !reflect.DeepEqual(test.output, got) {
			fmt.Printf("%+v, got: %v\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func twoSum(numbers []int, target int) []int {
	i, j := 0, len(numbers)-1
	for i < j {
		if numbers[i]+numbers[j] < target {
			i++
		} else if numbers[i]+numbers[j] > target {
			j--
		} else {
			return []int{i + 1, j + 1}
		}
	}
	return nil
}

// slowest
// func twoSum(numbers []int, target int) []int {
// 	for i := 0; i < len(numbers)-1; i++ {
// 		for j := len(numbers) - 1; j > i; j-- {
// 			if numbers[i]+numbers[j] == target {
// 				return []int{i + 1, j + 1}
// 			}
// 		}
// 	}
// 	return nil
// }
