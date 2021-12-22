package main

import (
	"fmt"
	"reflect"
)

func main() {
	tests := []struct {
		s      []byte
		output []byte
	}{
		{[]byte{'h', 'e', 'l', 'l', 'o'}, []byte{'o', 'l', 'l', 'e', 'h'}},
		{[]byte{'H', 'a', 'n', 'n', 'a', 'h'}, []byte{'h', 'a', 'n', 'n', 'a', 'H'}},
	}
	for _, test := range tests {
		reverseString(test.s)
		if !reflect.DeepEqual(test.output, test.s) {
			fmt.Printf("%+v\n", test)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func reverseString(s []byte) {
	i, j := 0, len(s)-1
	for i < j {
		s[i], s[j] = s[j], s[i]
		i++
		j--
	}
}
