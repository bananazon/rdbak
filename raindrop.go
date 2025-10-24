package main

import (
	"os"

	"github.com/bananazon/rdbak/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
