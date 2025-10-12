package main

import (
	"os"

	"github.com/gdanko/rdbak/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
