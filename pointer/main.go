package main

import "fmt"

type People struct {
	Name string
	Age  int
}

func main() {
	jack := People{Name: "jack", Age: 18}
	tom := &People{Name: "tom", Age: 18}

	fmt.Println(jack)
	addAge0(jack)
	fmt.Println(jack)
	addAge1(&jack)
	fmt.Println(jack)

	fmt.Printf("tom:\t\t[pointerValue: %p,\tvalue: %v,\tselfAddr: %p]\n", tom, tom, &tom)
	newPeople0(tom)
	fmt.Printf("tom:\t\t[pointerValue: %p,\tvalue: %v,\tselfAddr: %p]\n", tom, tom, &tom)
	newPeople1(tom)
	fmt.Printf("tom:\t\t[pointerValue: %p,\tvalue: %v,\tselfAddr: %p]\n", tom, tom, &tom)

	// {jack 18}
	// {jack 18}
	// {jack 19}
	// 0xc00000c030: &{tom 18}
	// p       0xc00000c0a8: &{tom 999}
	// 0xc00000c030: &{tom 18}
	// p       0xc00000c030: &{tomm 888}
	// 0xc00000c030: &{tomm 888}
}

func addAge0(p People) {
	p.Age++
}

func addAge1(p *People) {
	p.Age++
}

func newPeople0(p *People) {
	fmt.Printf("before p:\t[pointerValue: %p,\tvalue: %v,\tselfAddr: %p]\n", p, p, &p)
	p = &People{Name: "tom", Age: 999}
	fmt.Printf("after p:\t[pointerValue: %p,\tvalue: %v,\tselfAddr: %p]\n", p, p, &p)
}

func newPeople1(p *People) {
	fmt.Printf("before p:\t[pointerValue: %p,\tvalue: %v,\tselfAddr: %p]\n", p, p, &p)
	pp := People{Name: "tomm", Age: 888}
	*p = pp
	fmt.Printf("after p:\t[pointerValue: %p,\tvalue: %v,\tselfAddr: %p]\n", p, p, &p)
}
