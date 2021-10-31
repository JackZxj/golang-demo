package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"

	"github.com/JackZxj/golang-demo/sequencefile/sequencefile"
)

func main() {
	RWR()
}

func RWR() {
	var s sequencefile.SeqFile
	os.RemoveAll("data")
	s.Init()
	defer s.Close()
	var wg sync.WaitGroup
	wg.Add(400)
	for i := 0; i < 400; i++ {
		r := rand.Float32()
		if r < 0.6 {
			go func(i int) {
				defer wg.Done()
				index, err := s.Write([]byte(fmt.Sprintf("hello world, i am %d. 这篇文章将介绍Golang并发编程中常用到一种编程模式:context。本文将从为什么需要context出发,深入了解context的实现原理,以及了解如何使用context。 为什么需要context 在并发程序...这篇文章将介绍Golang并发编程中常用到一种编程模式:context。本文将从为什么需要context出发,深入了解context的实现原理,以及了解如何使用context。 为什么需要context 在并发程序...这篇文章将介绍Golang并发编程中常用到一种编程模式:context。本文将从为什么需要context出发,深入了解context的实现原理,以及了解如何使用context。 为什么需要context 在并发程序...", i)))
				if err != nil {
					log.Fatalf("write: %v", err)
				}
				log.Printf("######write %d\n", index)
			}(i)
			continue
		}
		go func(i int64) {
			defer wg.Done()
			b, err := s.Read(i)
			if err != nil {
				if !errors.Is(err, sequencefile.ErrIndexNotFound) {
					log.Fatalf("read: %v", err)
				}
				fmt.Println("oversize:", i)
			}
			log.Printf("######read %d: %q\n", i, b)
		}(int64(i) / 10)
	}
	wg.Wait()
}

func RWRR() {
	var s sequencefile.SeqFile
	os.RemoveAll("data")
	s.Init()
	defer s.Close()
	var wg sync.WaitGroup
	wg.Add(200)
	for i := 0; i < 200; i++ {
		go func(i int) {
			defer wg.Done()
			index, err := s.Write([]byte(fmt.Sprintf("hello world, i am %d", i)))
			if err != nil {
				log.Fatalf("write: %v", err)
			}
			log.Printf("######write %d\n", index)
		}(i)
	}
	wg.Wait()
	wg.Add(200)
	for i := int64(0); i < 200; i++ {
		go func(i int64) {
			defer wg.Done()
			b, err := s.Read(i)
			if err != nil {
				log.Fatalf("read: %v", err)
			}
			log.Printf("######read %d: %q\n", i, b)
		}(i)
	}
	wg.Wait()
}

// Sequential read and write and random read and write

func SWSR() {
	var s sequencefile.SeqFile
	os.RemoveAll("data")
	s.Init()
	for i := 0; i < 200; i++ {
		index, err := s.Write([]byte(fmt.Sprintf("hello world, i am %d", i)))
		if err != nil {
			log.Fatalf("write: %v", err)
		}
		log.Printf("######write %d\n", index)
	}
	for i := int64(0); i < 200; i++ {
		b, err := s.Read(i)
		if err != nil {
			log.Fatalf("read: %v", err)
		}
		log.Printf("######read: %q\n", b)
	}
	s.Close()
}
