package main

import "fmt"

func main() {
	a := []byte{1, 2, 3, 4}
	b := a[1:3]
	b[1] = 5
	fmt.Println(a, b)
	// --------
	b = append(b, 6)
	fmt.Println(a, b)
	// --------
	b = append(b, 7, 8) // 如果扩容超过原本长度就变成复制
	fmt.Println(a, b)
	fmt.Println("########################")

	c := []byte{1, 2, 3}
	fmt.Printf("%p\n", &c[1])
	c = c[1:2]
	fmt.Printf("%p\n", &c[0])
	c = []byte{c[0]}
	fmt.Printf("%p\n", &c[0])
	fmt.Println("########################")

	d := []byte{1, 2, 3}
	splitA(d)
	fmt.Printf("%v, %p\n", d, &d[0])
	splitB(&d)
	fmt.Printf("%v, %p\n", d, &d[0])
	fmt.Println("########################")

	var e [4]byte
	f := e[1:3]
	f[0], f[1] = 2, 3
	fmt.Println(e, f)
}

func splitA(a []byte) {
	a = a[:1]
}

func splitB(b *[]byte) {
	*b = (*b)[:1]
}
