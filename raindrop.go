package main

import (
	"os"

	"github.com/bananazon/raindrop/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
