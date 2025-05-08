package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var lock sync.RWMutex
var test string

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			demo2(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func demo(i int) {
	lock.RLock()
	defer lock.RUnlock()

	if test == "" {
		lock.RUnlock()
		defer lock.RLock()
		lock.Lock()
		defer lock.Unlock()
		if test == "" {
			time.Sleep(time.Duration(rand.Intn(10)) * time.Microsecond)
			test = "aaaa"
			fmt.Println(i, "changed")
		}
	}
	fmt.Println(i, test)
}

func demo2(i int) {
	if !isInit() {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Microsecond)
		if !isInit() {
			write("aaaa")
			fmt.Println(i, "changed")
		}
	}
	fmt.Println(i, read())
}

func read() string {
	lock.RLock()
	defer lock.RUnlock()

	return test
}

func isInit() bool {
	lock.Lock()
	defer lock.Unlock()

	return test != ""
}

func write(v string) {
	lock.Lock()
	defer lock.Unlock()

	test = v
}
