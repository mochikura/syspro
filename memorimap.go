package main

import (
	"fmt"
	"reflect"
)

var data0 int
var data1 = 10

func main() {
	c := 10
	// environ is inaccessible in Go
	// argv is inaccessible in Go
	fmt.Printf("stack:¥t¥t%p¥n", &c)
	fmt.Printf("bss:¥t¥t%p¥n", &data0)
	fmt.Printf("data:¥t¥t%p¥n", &data1)
	var ptr uintptr = reflect.ValueOf(main).Pointer()
	fmt.Printf("text:¥t¥t0x%x¥n", ptr)
}
