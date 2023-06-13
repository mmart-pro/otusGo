package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	var str = "Hello, OTUS!"
	str = stringutil.Reverse(str)
	fmt.Print(str)
}
