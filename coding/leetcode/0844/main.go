package main

import (
	"fmt"
)

func main() {
	tests := []struct {
		s    string
		t    string
		want bool
	}{
		{"", "", true},
		{"", "######", true},
		{"ab#c", "ad#c", true},
		{"ab##", "c#d#", true},
		{"a##c", "#a#c", true},
		{"a#c", "b", false},
		{"aa#c", "acb", false},
	}
	for _, test := range tests {
		if backspaceCompare(test.s, test.t) != test.want {
			fmt.Printf("%+v,got: %v\n", test, !test.want)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func backspaceCompare(s string, t string) bool {
	ss, ts := []byte(s), []byte(t)
	return string(ss[:after(ss)]) == string(ts[:after(ts)])
}

func after(ss []byte) int {
	index := 0
	for _, c := range ss {
		if c != '#' {
			ss[index] = c
			index++
		} else if index > 0 {
			index--
		}
	}
	return index
}
