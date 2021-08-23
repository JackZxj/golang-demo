package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	s := strings.Split(scanner.Text(), " ")
	// fmt.Println(s)
	if s[0] != "true" && s[0] != "false" {
		// fmt.Println(s[0], 0)
		fmt.Println("error")
		return
	}
	res := []string{s[0]}
	index := 0
	for i := 1; i < len(s); i++ {
		if index%2 == 0 {
			if s[i] == "true" || s[i] == "false" {
				// fmt.Println(s[i], i, 26)
				fmt.Println("error")
				return
			}
			res = append(res, s[i])
			index++
			continue
		}
		if s[i] != "true" && s[i] != "false" {
			// fmt.Println(s[i], i, 35)
			fmt.Println("error")
			return
		}
		if res[index] == "or" {
			res = append(res, s[i])
			index++
		} else {
			if s[i] != res[index-1] || s[i] == "false" {
				res[index-1] = "false"
			}
			res = res[0: index]
			index = index - 1
		}
		// fmt.Println(res)
	}
	if res[index] == "and" || res[index] == "or" {
		// fmt.Println(s[index], index, 50)
		fmt.Println("error")
		return
	}
	for _,v := range res {
		if v == "true" {
			fmt.Println(v)
			return
		}
	}
	fmt.Println("false")
}
