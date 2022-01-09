package main

import (
	"fmt"
	"reflect"
)

func main() {
	tests := []struct {
		n, k   int
		output [][]int
	}{
		{4, 3, [][]int{{1, 2, 3}, {1, 2, 4}, {1, 3, 4}, {2, 3, 4}}},
		{4, 2, [][]int{{1, 2}, {1, 3}, {1, 4}, {2, 3}, {2, 4}, {3, 4}}},
		{1, 1, [][]int{{1}}},
	}
	for _, test := range tests {
		if got := combine(test.n, test.k); !reflect.DeepEqual(got, test.output) {
			fmt.Printf("input %+v, got: %v\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

var res [][]int

func combine(n int, k int) [][]int {
	res = make([][]int, 0)
	forCombine([]int{}, n, k, 1)
	return res
}

func forCombine(base []int, n, k, now int) {
	// fmt.Printf("base: %v, n: %d, k:%d, now:%d\n", base, n, k, now)
	if k == 0 {
		// fmt.Println("append", base)
		tmp := make([]int, len(base))
		copy(tmp, base)
		res = append(res, tmp)
		return
	}
	for ; now <= n-k+1; now++ {
		base = append(base, now)
		forCombine(base, n, k-1, now+1)
		base = base[:len(base)-1]
	}
}
