package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type deployment struct {
	Force  []string `json:"force"`
	Kustom []string `json:"kustom"`
}

type uniqueDir struct {
	name []string
}

var outputJson []deployment

func main() {
	if len(os.Args) > 1 {
		var arr []string
		var uniqueDirectories uniqueDir

		dataJson := os.Args[1]

		_ = json.Unmarshal([]byte(dataJson), &arr)

		for _, path := range arr {
			directory := filepath.Dir(path)
			uniqueDirectories.name = append(uniqueDirectories.name, directory)
		}

		// Get only unique paths
		uniqueDirectories.name = unique(uniqueDirectories)

		// Setup Json Structure
		uniqueJson, err := uniqueDirectories.setupJson()
		//fmt.Printf("Unmarshaled Json: %v\n", uniqueJson)
		if err != nil {
			log.Fatal(err)
		}

		output, _ := json.Marshal(uniqueJson)
		fmt.Printf("Marshalled Json: %v", string(output))
	}
}

func unique(intSlice uniqueDir) []string {
	keys := make(map[string]string)
	var list []string
	for _, entry := range intSlice.name {
		if _, value := keys[entry]; !value {
			keys[entry] = entry
			list = append(list, entry)
		}
	}
	return list
}

// Json Setup
func (paths uniqueDir) setupJson() ([]deployment, error) {
	var force []string
	var kustom []string
	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for _, path := range paths.name {
		kustomizeFile := filepath.Join(workDir, path, "kustomization.yaml")

		if _, err = os.Stat(kustomizeFile); err == nil {
			kustom = append(kustom, "-k ./" + path)
		} else {

			force = append(force, "-f ./" + path)
		}
	}
	outputJson = append(outputJson, deployment{Force: force, Kustom: kustom})
	return outputJson, nil
}

//func (jsonValue deployment) testJson() (string, error) {
// if outputJson
//}