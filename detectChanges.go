package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// test string: ["Deployments/adguard/020-Issuer.yaml","Deployments/hello-world/010-Namespace.yaml"]
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
				fmt.Println("Found No Entry")
				directories[directory] = directory
				uniqueDirectories = append(uniqueDirectories, directory)
			}
		}
		// Format: ["values","values"]
		output := `['` + strings.Join(uniqueDirectories, `','`) + `']`
		fmt.Printf("%v", output)
	}
}
