package main

import "fmt"

func main() {
	a := []byte{1, 2, 3, 4}
	b := a[1:3]
	b[1] = 5
	fmt.Println("a         b           # 如果扩容超过原本长度就变成复制")
	fmt.Printf("%v %v       # &a[1]=%p,&b[0]=%p\n", a, b, &a[1], &b[0])
	// --------
	b = append(b, 6)
	fmt.Printf("%v %v     # &a[1]=%p,&b[0]=%p\n", a, b, &a[1], &b[0])
	// --------
	b = append(b, 7, 8) // 如果扩容超过原本长度就变成复制
	fmt.Printf("%v %v # &a[1]=%p,&b[0]=%p\n", a, b, &a[1], &b[0])
	fmt.Println("########################")

	c := []byte{1, 2, 3}
	fmt.Printf("%p\n", &c[1])
	c = c[1:2]
	fmt.Printf("%p\n", &c[0]) // 剪切是原数组，因此地址一致
	c = []byte{c[0]}
	fmt.Printf("%p\n", &c[0]) // 新建数组是复制
	fmt.Println("########################")

	d := []byte{1, 2, 3}
	fmt.Printf("修改前: %v, %p\n", d, &d[0])
	splitA(d)
	fmt.Printf("splitA: %v, %p\n\n", d, &d[0])

	d = []byte{1, 2, 3}
	fmt.Printf("修改前: %v, %p\n", d, &d[0])
	splitB(&d)
	fmt.Printf("splitB: %v, %p\n\n", d, &d[0])

	d = []byte{1, 2, 3}
	fmt.Printf("修改前: %v, %p\n", d, &d[0])
	splitC(d)
	fmt.Printf("splitC: %v, %p\n\n", d, &d[0])

	d = []byte{1, 2, 3}
	fmt.Printf("修改前: %v, %p\n", d, &d[0])
	splitD(d)
	fmt.Printf("splitD: %v, %p\n\n", d, &d[0])

	d = []byte{1, 2, 3}
	fmt.Printf("修改前: %v, %p\n", d, &d[0])
	splitE(d)
	fmt.Printf("splitE: %v, %p\n\n", d, &d[0])
	fmt.Println("########################")

	var e [4]byte
	f := e[1:3]
	f[0], f[1] = 2, 3
	fmt.Println(e, f) // slice 引用数组值，因此会修改 slice 会影响到原数组
	fmt.Println("########################")

	g := []byte{1, 2, 3, 4}
	h := []byte{5, 6}
	fmt.Println("修改前: ", g, h)
	copy(g[1:], h) // copy 作用在原 slice 上，不是引用
	h[0], h[1] = h[1], h[0]
	fmt.Println("修改前: ", g, h)
}

// 修改失败
func splitA(a []byte) {
	a = a[:1]
}

// 修改成功
func splitB(b *[]byte) {
	*b = (*b)[:1]
}

// 原地修改成功
func splitC(c []byte) {
	a := append(c[1:], c[:1]...)
	copy(c, a)
}

// 修改成功, copy 只能修改长度为 min(len(src), len(dst)) 的部分
func splitD(d []byte) {
	a := append(d[1:], d[:1]...)
	b := a[:1]
	copy(d, b)
}

// 修改成功, copy 只能修改长度为 min(len(src), len(dst)) 的部分
func splitE(d []byte) {
	a := append(d[1:], d...)
	copy(d, a)
}
