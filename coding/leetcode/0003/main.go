package main

import "fmt"

func main() {
	tests := []struct {
		input  string
		output int
	}{
		{"abcabcbb", 3},
		{"bbbbb", 1},
		{"pwwkew", 3},
		{"", 0},
		{"au", 2},
		{"abb", 2},
		{"aab", 2},
	}
	for _, test := range tests {
		if got := lengthOfLongestSubstring(test.input); got != test.output {
			fmt.Printf("%+v, got: %d\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func lengthOfLongestSubstring(s string) int {
	if len(s) < 2 {
		return len(s)
	}
	ss := []byte(s)
	low, high := 0, 1
	max := 0
	for ; high < len(ss); high++ {
		for i := low; i < high; i++ {
			if ss[high] == ss[i] {
				if high-low > max {
					max = high - low
				}
				low = i + 1
				break
			}
		}
	}
	if high-low > max {
		return high - low
	}
	return max
}
