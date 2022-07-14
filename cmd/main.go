package main

import (
	"fmt"
	"os"

	"github.com/mx-psi/grater/internal"
)

func main() {
	cmd := internal.NewCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
