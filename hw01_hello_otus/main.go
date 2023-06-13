package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	// Place your code here.
	var str string = "Hello, OTUS!"
	str = stringutil.Reverse(str)
	fmt.Print(str)
}
