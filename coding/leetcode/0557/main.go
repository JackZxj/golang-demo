package main

import (
	"fmt"
)

func main() {
	tests := []struct {
		s      string
		output string
	}{
		{"", ""},
		{"a b", "a b"},
		{"abc", "cba"},
		{"ab  c", "ba  c"},
		{"Let's take LeetCode contest", "s'teL ekat edoCteeL tsetnoc"},
		{" Let's take LeetCode contest", " s'teL ekat edoCteeL tsetnoc"},
	}
	for _, test := range tests {
		if got := reverseWords(test.s); got != test.output {
			fmt.Printf("%+v, got: %q\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func reverseWords(s string) string {
	i := 0
	ss := []byte(s)
	for j := range ss {
		if ss[j] == ' ' || j == len(ss)-1 {
			t := j - 1
			if j == len(ss)-1 {
				t = j
			}
			for i < t {
				ss[i], ss[t] = ss[t], ss[i]
				i++
				t--
			}
			i = j + 1
		}
	}
	return string(ss)
}
