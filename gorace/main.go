package main

import (
	"fmt"
	"sync"
)

// race 不会 panic
func main() {
	var race int
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if e := recover(); e != nil {
					fmt.Println(e)
				}
			}()
			race++
			fmt.Println(race)
		}()
	}
	wg.Wait()
}
