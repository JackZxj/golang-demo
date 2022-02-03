package main

import (
	"fmt"
	"reflect"
)

func main() {
	tests := []struct {
		s    string
		p    string
		want []int
	}{
		{"a", "b", []int{}},
		{"aaa", "aab", []int{}},
		{"aaa", "a", []int{0, 1, 2}},
		{"aaab", "a", []int{0, 1, 2}},
		{"aaabaa", "a", []int{0, 1, 2, 4, 5}},
		{"abab", "ab", []int{0, 1, 2}},
		{"ababcaabbaccccc", "ab", []int{0, 1, 2, 6, 8}},
		{"cbaebabacd", "abc", []int{0, 6}},
		{"abaacbabc", "abc", []int{3, 4, 6}},
		{"ababababbacbabc", "abc", []int{8, 9, 10, 12}},
		{"aabbccbbcaaa", "bbcc", []int{2, 3, 4, 5}},
	}
	for _, test := range tests {
		got := findAnagrams(test.s, test.p)
		if !reflect.DeepEqual(got, test.want) {
			fmt.Printf("%+v, got: %v\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func findAnagrams(s string, p string) []int {
	var (
		result = make([]int, 0, len(s))
		dist   = make(map[byte]int, len(p))
		start  int
	)
	if len(s) < len(p) {
		return result
	}
	for i := range p {
		dist[p[i]]++
	}
	for i := range s {
		if _, exist := dist[s[i]]; exist {
			if dist[s[i]] > 0 {
				dist[s[i]]--
				if (i - start + 1) != len(p) {
					continue
				}
				result = append(result, start)
				dist[s[start]]++
				start++
				continue
			}
		}
		for ; s[start] != s[i]; start++ {
			if _, exist := dist[s[start]]; exist {
				dist[s[start]]++
			}
		}
		start++
	}
	return result
}
