package main

import (
	"github.com/spf13/cobra"
)

func main() {
	var rootCommand = &cobra.Command{
		Use:   "malguem",
		Short: "Boilerplate code generator",
		Long:  "A CLI tool to generate boilerplate code for multiple languages using templates.",
	}

	rootCommand.Execute()
}
