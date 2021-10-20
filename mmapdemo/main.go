package main

import (
	"fmt"
	"io/ioutil"
	"log"
	// "net/http"
	// _ "net/http/pprof"
	"os"
	"syscall"
)

// so...
// mmap is more faster when reading big a file,
// and using more less memory
func main() {
	f, err := os.Open("/root/board-offline-installer-7.4.tgz")
	if err != nil {
		log.Fatalf("open file: %v", err)
	}
	defer f.Close()
	var data []byte
	if len(os.Args) < 2 {
		fmt.Println("miss parameter 'm' or 'i'")
		return
	}
	flag := os.Args[1]
	switch flag {
	case "m":
		data = mmapRead(f)
	case "i":
		data = osRead(f)
	default:
		fmt.Println("accept 'm' or 'i' only")
		return
	}
	fmt.Println(data[len(data)/2])

	// log.Println(http.ListenAndServe("0.0.0.0:10000", nil))
}

func osRead(f *os.File) []byte {
	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("read file: %v", err)
	}
	return data
}

func mmapRead(f *os.File) []byte {
	fstat, _ := f.Stat()
	b, err := syscall.Mmap(int(f.Fd()), 0, int(fstat.Size()), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	return b
}

/*
[root@10 mmapdemo]# ls -lh /root/board-offline-installer-7.4.tgz
-rw-r--r-- 1 root root 1.9G Dec 21  2020 /root/board-offline-installer-7.4.tgz

[root@10 mmapdemo]# go build -o mmapdemo

[root@10 mmapdemo]# time ./mmapdemo m
171

real    0m0.012s
user    0m0.002s
sys     0m0.002s
[root@10 mmapdemo]# time ./mmapdemo m
171

real    0m0.003s
user    0m0.001s
sys     0m0.001s
[root@10 mmapdemo]# time ./mmapdemo m
171

real    0m0.002s
user    0m0.001s
sys     0m0.002s
[root@10 mmapdemo]# time ./mmapdemo i
171

real    0m28.898s
user    0m8.006s
sys     0m4.750s
[root@10 mmapdemo]# time ./mmapdemo i
171

real    0m28.705s
user    0m7.113s
sys     0m5.002s
[root@10 mmapdemo]# time ./mmapdemo i
171

real    0m34.528s
user    0m6.237s
sys     0m4.229s
*/
