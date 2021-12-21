package main

import "fmt"

var badVersion int

func main() {
	tests := []struct {
		input int
		bad   int
	}{
		{1, 1},
		{4, 1},
		{4, 2},
		{4, 3},
		{4, 4},
		{5, 1},
		{5, 2},
		{5, 3},
		{5, 4},
		{5, 5},
		{7, 3},
		{7, 7},
	}
	for _, test := range tests {
		badVersion = test.bad
		if got := firstBadVersion(test.input); got != test.bad {
			fmt.Println(test.input, test.bad, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

/**
 * Forward declaration of isBadVersion API.
 * @param   version   your guess about first bad version
 * @return 	 	      true if current version is bad
 *			          false if current version is good
 * func isBadVersion(version int) bool;
 */
func isBadVersion(version int) bool {
	return version >= badVersion
}

func firstBadVersion(n int) int {
	// return forSearch(n)
	return bSearch(1, n)
}

func forSearch(n int) int {
	i := 1
	for i < n {
		t := (i + n) / 2
		if isBadVersion(t) {
			n = t
		} else {
			i = t + 1
		}
	}
	return i
}

func bSearch(l, h int) int {
	if l >= h {
		return l
	}
	if t := (l + h) / 2; isBadVersion(t) {
		return bSearch(l, t)
	} else {
		return bSearch(t+1, h)
	}
}
