package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) > 1 {
		var arr []string
		var uniqueDirectories []string

		dataJson := os.Args[1]
		directories := make(map[string]string)

		_ = json.Unmarshal([]byte(dataJson), &arr)

		for _, path := range arr {
			directory := filepath.Dir(path)
			if _, ok := directories[directory]; !ok {
				directories[directory] = directory
				uniqueDirectories = append(uniqueDirectories, directory)
			}
		}

		output, _ := json.Marshal(&uniqueDirectories)
		fmt.Printf("%v", output)
	}
}
