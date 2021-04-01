package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		directories := make(map[int]string)
		for i, arg := range os.Args {
			if i > 0 {
				directories[i] = arg
			}
		}
		fmt.Println(directories)
	}
}
