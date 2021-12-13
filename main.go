package main

import (
	"fmt"
	"os"

	"nomad-gitops-operator/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
