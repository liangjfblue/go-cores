package main

import (
	"cgo-demo/service"
	"fmt"
)

func main() {
	a := "1"
	b := "2"

	fmt.Println(service.Append(a, b))
}
