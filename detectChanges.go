package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func unique(intSlice []string) []string {
	keys := make(map[string]string)
	var list []string
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = entry
			list = append(list, entry)
		}
	}
	return list
}

func main() {
	if len(os.Args) > 1 {
		var arr []string
		var uniqueDirectories []string

		dataJson := os.Args[1]

		_ = json.Unmarshal([]byte(dataJson), &arr)

		for _, path := range arr {
			directory := filepath.Dir(path)
			uniqueDirectories = append(uniqueDirectories, directory)
		}

		output, _ := json.Marshal(unique(uniqueDirectories))
		fmt.Printf("%v", string(output))
	}
}
