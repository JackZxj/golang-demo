package main

import (
	"fmt"
	"sync"
)

func main() {
	// aa()
	bb()
}

// map 随机读写删除需要加锁
func aa() {
	var mm = map[string]string{
		"a": "a",
		"b": "b",
		"c": "c",
		"d": "d",
		"e": "e",
		"f": "f",
		"g": "g",
		"h": "h",
		"i": "i",
		"j": "j",
		"k": "k",
		"l": "l",
		"m": "m",
		"n": "n",
		"o": "o",
		"p": "p",
		"q": "q",
		"r": "r",
		"s": "s",
		"t": "t",
		"u": "u",
		"v": "v",
		"w": "w",
		"x": "x",
		"y": "y",
		"z": "z",
	}
	var wg sync.WaitGroup
	var rw sync.RWMutex
	var l = []string{"bb", "cc", "dd", "ee", "ff", "gg", "hh", "ii", "aaa", "b", "c", "d", "ew", "f", "g", "hw", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	for ii := range l {
		wg.Add(1)
		i := ii
		go func() {
			defer wg.Done()
			rw.RLock()
			_, exist := mm[l[i]]
			rw.RUnlock()
			if exist {
				rw.Lock()
				delete(mm, l[i])
				rw.Unlock()
				fmt.Println("exists", l[i])
				return
			}
			fmt.Println("not exists", l[i])
		}()
	}
	wg.Wait()
	fmt.Println("done1")
	for kk, vv := range mm {
		k, v := kk, vv
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("left", k, v)
		}()
	}
	wg.Wait()
	fmt.Println("done2")
}

func bb() {
	m := make(map[string]string)
	editmap(m)
	fmt.Printf("修改成功: %v\n", m)
	m["hello"] = "aaa"
	fmt.Printf("修改成功: %v\n", m)
	editmap(m)
	fmt.Printf("修改成功: %v\n", m)
	editmap2(m)
	fmt.Printf("修改失败: %v\n", m)
	editmap3(&m)
	fmt.Printf("修改失败: %v\n", m)
	editmap4(&m)
	fmt.Printf("修改成功: %v\n", m)
}

// 修改成功
func editmap(m map[string]string) {
	m["hello"] = "editmap"
}

// 不能替换
func editmap2(m map[string]string) {
	mm := make(map[string]string)
	mm["gg"] = "editmap2"
	m = mm
}

// 不能替换
func editmap3(m *map[string]string) {
	mm := make(map[string]string)
	mm["gg"] = "editmap3"
	m = &mm
}

// 可以替换
func editmap4(m *map[string]string) {
	mm := make(map[string]string)
	mm["gg"] = "editmap4"
	*m = mm
}
