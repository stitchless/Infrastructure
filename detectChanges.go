package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)
// test string: ["Deployments/adguard/010-Namespace.yaml","Deployments/hello-world/010-Namespace.yaml","Deployments/hello-world/010-Namespace.yaml"]
func main() {
	if len(os.Args) > 1 {
		//directories := make(map[string]string)
		var changedDirectories []string

		err := json.Unmarshal([]byte(os.Args[1]), &changedDirectories)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(changedDirectories)

		//for i, arg := range os.Args {
		//	if i > 0 {
		//		directory := filepath.Dir(arg)
		//		if _, ok := directories[directory]; !ok {
		//			fmt.Println("Found No Entry")
		//			directories[directory] = directory
		//			changedDirectories = append(changedDirectories, directory)
		//		}
		//		fmt.Println("Found Entry for: " + directory)
		//	}
		//}
		// Format: ["values","values"]
		//output := `['` + strings.Join(changedDirectories, `','`) + `']`
		fmt.Println(changedDirectories)
	}
}
