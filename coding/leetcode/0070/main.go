package main

import "fmt"

func main() {
	tests := []struct {
		input, output int
	}{
		{1, 1},
		{2, 2},
		{3, 3},
		{44, 1134903170},
	}
	for _, test := range tests {
		if got := climbStairs(test.input); got != test.output {
			fmt.Printf("input: %+v, got: %v\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

// 斐波那契数列
func climbStairs(n int) int {
	j, k := 1, 1
	for i := 2; i <= n; i++ {
		j, k = k, j+k
	}
	return k
}

// ####################### solution 3 ##########################

// var cache []int

// func climbStairs(n int) int {
// 	cache = make([]int, n+1)
// 	cache[0] = 1
// 	cache[1] = 1
// 	return countN(n)
// }

// func countN(n int) int {
// 	if cache[n] == 0 {
// 		cache[n] = countN(n-1) + countN(n-2)
// 	}
// 	return cache[n]
// }

// ####################### solution 1 ##########################

// // timeout
// func climbStairs(n int) int {
// 	if n < 3 {
// 		return n
// 	}
// 	return climbStairs(n-1) + climbStairs(n-2)
// }
