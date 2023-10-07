package main

import (
	"fmt"
	"os"
)

func main() {
	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	os.Exit(RunCmd(os.Args[2:], env))
}
