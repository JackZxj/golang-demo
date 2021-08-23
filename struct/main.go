package main

import "fmt"

type Bar interface {
	setVal(newVal int)
	getVal() int
}

type Ball struct {
	val int
}

func (b Ball) getVal() int {
	return b.val
}
func (b *Ball) setVal(newVal int) {
	b.val = newVal
}
func (b Ball) setV(newVal int) {
	b.val = newVal
}

type Person1 struct {
	person1 *Ball
}
type Person2 struct {
	person2 Ball
}

func main() {
	ball := Ball{val: 101}
	fmt.Println(ball)

	// var person1 Person1 // 初始化为{person1: nil}

	// person1 := Person1{&Ball{val: 10}} // 
	// person1.person1 = &ball
	person1 := &Person1{&Ball{val: 10}} // 等效于下面这行
	// person1 := &Person1{&ball}
	fmt.Println(person1.person1.getVal())
	person1.person1.setVal(1000)
	fmt.Println(person1.person1.getVal()) // changed the value
	person1.person1.setV(111)
	fmt.Println(person1.person1.getVal()) // did not change the value
	
	
	fmt.Println()
	var person2 Person2 // 初始化为{person2: Ball{val： 0}}

	// person2 := Person2{Ball{val: 20}}
	person2.person2 = ball
	fmt.Println(person2.person2.getVal())
	person2.person2.setVal(2000)
	fmt.Println(person2.person2.getVal())
}
