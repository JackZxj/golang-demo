package main

import "fmt"

func main() {
	tests := []struct {
		n      int
		output bool
	}{
		{-1, false},
		{0, false},
		{1, true},
		{2, true},
		{3, false},
		{5, false},
		{16, true},
	}
	for _, test := range tests {
		if isPowerOfTwo(test.n) != test.output {
			fmt.Printf("%+v\n", test)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func isPowerOfTwo(n int) bool {
	if n == 0 {
		return false
	}
	return n&(n-1) == 0
}

func isPowerOfTwo2(n int) bool {
	t := 1
	for n > t {
		t = t << 1
	}
	return n == t
}
