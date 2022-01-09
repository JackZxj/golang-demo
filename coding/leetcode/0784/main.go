package main

import (
	"fmt"
	"reflect"
)

func main() {
	tests := []struct {
		s      string
		output []string
	}{
		{"a1b2", []string{"a1b2", "a1B2", "A1b2", "A1B2"}},
		{"3z4", []string{"3z4", "3Z4"}},
		{"12345", []string{"12345"}},
	}
	for _, test := range tests {
		got := letterCasePermutation(test.s)
		if !reflect.DeepEqual(got, test.output) {
			fmt.Printf("input: %+v, got: %v\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

var res []string

func letterCasePermutation(s string) []string {
	res = []string{}
	forL([]byte(s), -1)
	return res
}

func forL(base []byte, cur int) {
	next := cur
	existNext := false
	for i := cur + 1; i < len(base); i++ {
		if base[i] > '9' {
			next = i
			existNext = true
			break
		}
	}
	if !existNext {
		res = append(res, string(base))
		return
	}
	if base[next] < 'a' {
		// base[next]是大写，转小写
		base[next] += 32
	}
	forL(base, next)
	base[next] -= 32
	forL(base, next)
}
