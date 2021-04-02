package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	//fmt.Println("Reading Environment Variables...")
	//var inputDirectories string
	//inputDirectories = os.Getenv("TEST")
	//fmt.Printf("Found Directories: %s\n", inputDirectories)


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
		}
		fmt.Print(output)
	}
}