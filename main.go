package main

import (
	"fmt"
	"os"

	"lab_cli/cmd"
)

func main() {
	if err := cmd.Execute(os.Args[1:]); err != nil {
		fmt.Println("Unable to run command: ", err)
		os.Exit(1)
	}
}
