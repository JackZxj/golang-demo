package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	Name string
	Age  int
}

type Teacher struct {
	Name string
	Age  int
}

type SuperStudent struct {
	Name  string
	Age   int
	Title map[string]string // 含有 map, slice 的对象无法直接比较
	Score []int
}

func main() {
	s1 := Student{Name: "s1", Age: 12}
	s2 := Student{Name: "s1", Age: 12}
	fmt.Println("s1 == s2: ", s1 == s2)     // true
	fmt.Println("&s1 == &s2: ", &s1 == &s2) // false

	t1 := Teacher{Name: "s1", Age: 12}
	// fmt.Println("t1 == s1: ", t1 == s1) // 无法跨类型比较
	t1AsStudent := Student(t1)
	fmt.Println("t1AsStudent == s1: ", t1AsStudent == s1)     // true
	fmt.Println("&t1AsStudent == &s1: ", &t1AsStudent == &s1) // false

	ss1 := SuperStudent{Name: "ss1", Age: 12, Title: make(map[string]string), Score: []int{99, 99, 100}}
	ss2 := SuperStudent{Name: "ss1", Age: 12, Title: make(map[string]string), Score: []int{99, 99, 100}}
	// fmt.Println("ss1 == ss2: ", ss1 == ss2) // 结构体无法直接比较
	fmt.Println("ss1 == ss2: ", reflect.DeepEqual(ss1, ss2)) // true
	fmt.Println("&ss1 == &ss2: ", &ss1 == &ss2)              // false
}
