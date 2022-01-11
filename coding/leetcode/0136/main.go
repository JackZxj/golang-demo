package main

import "fmt"

func main() {
	tests := []struct {
		in  []int
		out int
	}{
		{[]int{2, 2, 1}, 1},
		{[]int{4, 1, 2, 1, 2}, 4},
	}
	for _, test := range tests {
		if got := singleNumber(test.in); got != test.out {
			fmt.Printf("%+v,got: %d\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func singleNumber(nums []int) (out int) {
	for i := range nums {
		out ^= nums[i]
	}
	return
}
