package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) > 1 {
		directories := make(map[string]string)
		var output []string

		for i, arg := range os.Args {
			if i > 0 {
				directory := filepath.Dir(arg)
				if _, ok := directories[directory]; !ok {
					directories[directory] = directory
					output = append(output, directory)
				}
			}
			fmt.Println(output)
		}
	}
}