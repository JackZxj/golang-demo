package main

import "fmt"

import "github.com/JackZxj/golang-demo/getpwd/pwd"

func main() {
	fmt.Println("###read go.mod:\n", pwd.Run())
}
