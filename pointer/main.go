package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

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
	testClient4()
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

type CLient interface {
	Get(url string) (resp *http.Response, err error)
}

func NewClient() CLient {
	c := http.Client{}
	fmt.Printf("cli new: %p\n", &c)
	return &c
}

// Pointer changes do not affect running functions
func testClient() {
	c := NewClient()
	done := make(chan int)
	go func() {
		time.Sleep(time.Second)
		fmt.Printf("cli before: %p\n", c)
		c = NewClient()
		fmt.Printf("cli changed: %p\n", c)
		done <- 1
	}()

	for i := 0; i < 3; i++ {
		fmt.Printf("for: %p\n", c)
		if _, err := c.Get("http://103.209.9.98"); err != nil {
			fmt.Println(i, "got err:", err)
		}
	}
	<-done
}

// Pointer wont change in new routine
func testClient2() {
	c := NewClient()
	done := make(chan int)

	getC := func() CLient {
		return c
	}

	go func() {
		time.Sleep(time.Second)
		fmt.Printf("cli before: %p\n", c)
		c = NewClient()
		fmt.Printf("cli changed: %p\n", c)
		done <- 1
	}()

	go func() {
		c := getC()
		for i := 0; i < 3; i++ {
			fmt.Printf("for: %p\n", c)
			if _, err := c.Get("http://103.209.9.98"); err != nil {
				fmt.Println(i, "got err:", err)
			}
		}
		done <- 1
	}()

	<-done
	<-done
}

// Pointer wont change in new routine, as same as testClient4
func testClient3() {
	c := NewClient()
	done := make(chan int)

	getC := func() CLient {
		return c
	}

	go func() {
		time.Sleep(time.Second)
		fmt.Printf("cli before: %p\n", c)
		c = NewClient()
		fmt.Printf("cli changed: %p\n", c)
		done <- 1
	}()

	go func() {
		c := getC()
		for i := 0; i < 3; i++ {
			fmt.Printf("for: %p\n", c)
			if _, err := c.Get("http://103.209.9.98"); err != nil {
				if strings.Contains(err.Error(), "network is unreachable") {
					c = getC()
				}
				fmt.Println(i, "got err:", err)
			}
		}
		done <- 1
	}()

	<-done
	<-done
}

// Pointer wont change in new routine, as same as testClient3
func testClient4() {
	c := NewClient()
	cc := &c
	done := make(chan int)

	getC := func() *CLient {
		return cc
	}

	go func() {
		time.Sleep(time.Second)
		fmt.Printf("cli before: cc:%p *cc:%p\n", cc, *cc)
		nc := NewClient()
		cc = &nc
		fmt.Printf("cli changed: cc:%p *cc:%p\n", cc, *cc)
		done <- 1
	}()

	go func() {
		c := getC()
		for i := 0; i < 3; i++ {
			fmt.Printf("for: c:%p *c:%p\n", c, *c)
			if _, err := (*c).Get("http://103.209.9.98"); err != nil {
				if strings.Contains(err.Error(), "network is unreachable") {
					c = getC()
				}
				fmt.Println(i, "got err:", err)
			}
		}
		done <- 1
	}()

	<-done
	<-done
}
