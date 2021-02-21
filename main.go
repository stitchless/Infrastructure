package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	// Allow user to supply custom working directory
	targetDirectory := "Deployments"
	if len(os.Args) > 1 && os.Args[1] != "" {
		targetDirectory = os.Args[1]
	}

	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	_ = filepath.Walk(filepath.Join(workingDir, targetDirectory), func(path string, info os.FileInfo, err error) error {
		if err == nil && info.Name() == "kustomization.yaml" {
			parseFile(path)
		}
		return nil
	})
}

func parseFile(path string) {
	str := regexp.MustCompile(`\${{2}.+(?:secrets\.(.+)).+}{2}`)
	file, err := os.ReadFile(path)
	check(err)
	fileStat, err := os.Stat(path)
	check(err)
	permissions := fileStat.Mode()

	isExist := str.MatchString(string(file))
	if isExist {
		newContents := str.ReplaceAllStringFunc(string(file), func(element string) string {
			envKey := str.FindStringSubmatch(element)[1]
			env, envBool := os.LookupEnv(envKey)
			if envBool {
				fmt.Printf("Found Environment Variable for : %s\n", envKey)
				return env
			}
			panic("Could not assign a valid value for env key: " + envKey + ".  Environment variable not found.")
		})
		err = os.WriteFile(path, []byte(newContents), permissions)
		if err != nil {
			panic(err)
		}
	}
}

// Error Checking
func check(e error) {
	if e != nil {
		panic(e)
	}
}
