//usr/bin/env go run $0 "$@"; exit
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	if len(os.Args) > 1 {
		app := "echo"
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
		fmt.Println(output)
		// Send the results to the outputs of the github script step via echo command
		formattedArg := fmt.Sprintf("::set-output name=changed_output::%+v", output)
		cmd := exec.Command(app, formattedArg)
		stdout, err := cmd.Output()

		if err != nil {
			println(err.Error())
		}
		print(string(stdout))
	}
}