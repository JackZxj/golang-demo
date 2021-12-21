package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	fileLink(true)
	res := make([]string, 0)
	err := walkPath("data", 2, &res)
	if err != nil {
		log.Fatal("rm datapath:", err)
	}
	log.Println(res)
}

// windows not supports Symlink, and supports Link partly:
// https://docs.microsoft.com/zh-cn/windows/win32/fileio/filesystem-functionality-comparison?redirectedfrom=MSDN
func fileLink(hard bool) {
	var err error
	datapath := "data"
	err = os.RemoveAll(datapath)
	if err != nil {
		log.Fatal("rm datapath:", err)
	}

	err = os.MkdirAll(datapath, 0755)
	if err != nil {
		log.Fatalf("create data path: %v", err)
	}
	origin := filepath.Join(datapath, "original.txt")
	snap1 := filepath.Join(datapath, "snap-1", "original.txt")
	snap2 := filepath.Join(datapath, "snap-2", "original.txt")
	err = ioutil.WriteFile(origin, []byte("hello,world!hello,world!hello,world!hello,world!hello,world!hello,world!hello,world!hello,world!hello,world!hello,world!"), 0600)
	if err != nil {
		log.Fatalln(err)
	}
	// make snap
	err = os.MkdirAll(filepath.Join(datapath, "snap-1", "son"), 0755)
	// ioutil.WriteFile(filepath.Join(datapath, "snap-1", "son", "son.txt"), []byte("sssssssssonsssssssssonsssssssssonsssssssssonsssssssssonsssssssssonssssssssson"), 0600)
	if err != nil {
		log.Fatalf("create data path: %v", err)
	}
	if hard {
		err = os.Link(origin, snap1)
	} else {
		err = os.Symlink(origin, snap1)
	}
	if err != nil {
		log.Fatalln(err)
	}
	s1, _ := os.Open(snap1)
	b1, _ := ioutil.ReadAll(s1)
	log.Printf("snap1: %s\n", b1)

	err = os.MkdirAll(filepath.Join(datapath, "snap-2"), 0755)
	if err != nil {
		log.Fatalf("create data path: %v", err)
	}
	if hard {
		err = os.Link(origin, snap2)
	} else {
		err = os.Symlink(origin, snap2)
	}
	if err != nil {
		log.Fatalln(err)
	}
	s2, _ := os.Open(snap2)

	b2, _ := ioutil.ReadAll(s2)
	log.Printf("snap2: %s\n", b2)

	// Link path not supports

	// pathsnap := filepath.Join(datapath, "pathsnap")
	// // os.MkdirAll(pathsnap, 0755)
	// err = os.Symlink(filepath.Join(datapath, "snap-1"), pathsnap)
	// err = os.Link(filepath.Join(datapath, "snap-1"), pathsnap)

	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// p1, _ := os.Open(filepath.Join(pathsnap, "son", "son.txt"))
	// b3, _ := ioutil.ReadAll(p1)
	// log.Printf("pathsnap: %s\n", b3)
}

func walkPath(path string, maxDepth int, result *[]string) error {
	fileInfos, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, fi := range fileInfos {
		name := filepath.Join(path, fi.Name())
		if !fi.IsDir() {
			*result = append(*result, name)
			log.Println(name)
		} else {
			if maxDepth > 1 {
				err = walkPath(name, maxDepth-1, result)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
