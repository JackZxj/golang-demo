package main

import "fmt"

func main() {
	tests := []struct {
		in, out uint32
	}{
		{0b00000010100101000001111010011100, 0b00111001011110000010100101000000},
		{0b11111111111111111111111111111101, 0b10111111111111111111111111111111},
		{0b00000000000000000000000000000001, 0b10000000000000000000000000000000},
	}
	for _, test := range tests {
		if got := reverseBits(test.in); got != test.out {
			fmt.Printf("%+v,got: %d\n", test, got)
			panic("not ok")
		}
	}
	fmt.Println("ok")
}

func reverseBits(num uint32) (out uint32) {
	for i := 0; i < 32; i++ {
		out = (out << 1) | (num >> i & 1)
	}
	return out
}
