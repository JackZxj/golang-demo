/*
# interface 类型的巧用

interface 是可比较类型，同其他基础类型一样可以用 == 判断，
也可以作为 map 的 key

参考：
 * https://zhuanlan.zhihu.com/p/522492698
 * https://www.veaxen.com/golang%E6%8E%A5%E5%8F%A3%E5%80%BC%EF%BC%88interface%EF%BC%89%E7%9A%84%E6%AF%94%E8%BE%83%E6%93%8D%E4%BD%9C%E5%88%86%E6%9E%90.html

*/
package main

import "fmt"

type peopleer interface {
	name() string
}

type chinese struct {
	fullName string
}

func (d *chinese) name() string {
	return d.fullName
}

func newChinese(name string) chinese {
	return chinese{fullName: name}
}

type english struct {
	firstName string
	lastName  string
}

func (d *english) name() string {
	return d.firstName + "·" + d.lastName
}

func newEnglish(name1, name2 string) english {
	return english{firstName: name1, lastName: name2}
}

var m map[peopleer]struct{}

func collection(ps ...peopleer) {
	if m == nil {
		m = make(map[peopleer]struct{})
	}
	for _, p := range ps {
		m[p] = struct{}{}
	}
}

func printCollection() {
	if m == nil {
		return
	}
	for k := range m {
		fmt.Println(k.name())
	}
}

func comparePeople(a, b peopleer) bool {
	return a == b
}

func main() {
	d1 := newChinese("张三")
	d2 := newEnglish("three", "zhang")
	d3 := chinese{fullName: "王五"}
	d4 := english{"five", "wang"}

	collection(&d1, &d2, &d3, &d4)
	printCollection()

	d5 := chinese{fullName: "张三"}
	fmt.Println("d5==d1?                   ", d5 == d1)
	fmt.Println("&d5==&d1?                 ", &d5 == &d1)
	fmt.Println("people(d5)==people(d1)?   ", comparePeople(&d1, &d5)) // 实际上是因为指针地址不同，道理通 &d5 == &d1。运行的时候会自动调用 runtime.convT2I(SB) 将 struct 转为 interface
}
