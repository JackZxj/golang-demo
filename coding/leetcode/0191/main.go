package main

import "fmt"

func main() {
	tests := []struct {
		n      uint32
		output int
	}{
		{0b00000000000000000000000000000000, 0},
		{0b00000000000000000000000000001011, 3},
		{0b00000000000000000000000010000000, 1},
		{0b11111111111111111111111111111101, 31},
		{0b11111111111111111111111111111111, 32},
	}
	for _, test := range tests {
		if got := hammingWeight(test.n); got != test.output {
			fmt.Printf("%+v,got: %d\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}
func hammingWeight(num uint32) int {
	count := 0
	for num != 0 {
		if num&1 == 1 {
			count++
		}
		num = num >> 1
	}
	return count
}
