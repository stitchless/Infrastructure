package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) > 1 {
		directories := make(map[string]string)
		var changedDirectories []string

		for i, arg := range os.Args {
			if i > 0 {
				directory := filepath.Dir(arg)
				if _, ok := directories[directory]; !ok {
					directories[directory] = directory
					changedDirectories = append(changedDirectories, directory)
				}
			}
		}
		//fmt.Println(changedDirectories)
		// Format: ["values","values"]
		output := `["` + strings.Join(changedDirectories, `","`) + `"]`
		fmt.Println(output)
	}
}
