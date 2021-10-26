package main

import (
	"bytes"
	"log"
	"os"

	m1 "github.com/edsrzf/mmap-go"
)

func main() {
	f, err := os.OpenFile("test.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	mmap1(f)
}

func mmap1(f *os.File) {
	if err := f.Truncate(1024); err != nil {
		log.Fatal(err)
	}

	mmap, err := m1.Map(f, m1.RDWR, 0)
	if err != nil {
		log.Fatalf("mmap error: %v", err)
	}
	defer mmap.Unmap()
	for index, bb := range []byte("Hello mmap1") {
		mmap[index] = bb
	}
	mmap.Flush()
	if err := f.Truncate(2048); err != nil {
		log.Fatal(err)
	}
	mmap, err = m1.Map(f, m1.RDWR, 0)
	if err != nil {
		log.Fatalf("mmap error: %v", err)
	}

	buf := bytes.NewBuffer(mmap)
	log.Printf("length: %v; ospage: %d\n", buf.Len(), os.Getpagesize())
	tmp := buf.Next(4096)
	log.Printf("%v: %d, %d\n", tmp, len(tmp), buf.Len())
}

// func mmap2(f *os.File) {
// 	if err := f.Truncate(8192); err != nil {
// 		log.Fatal(err)
// 	}
// 	mmap,err := m1.MapRegion()
// }