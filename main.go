package main

import (
	"github.com/spf13/cobra"
	"malguem/cmd"
)

var rootCommand = &cobra.Command{
	Use:   "malguem",
	Short: "Boilerplate code generator",
	Long:  "A CLI tool to generate boilerplate code for multiple languages using templates.",
}

func main() {
	rootCommand.AddCommand(cmd.InitCommand)
	rootCommand.AddCommand(cmd.CreateCommad)
	rootCommand.AddCommand(cmd.MakeCommand)
	rootCommand.AddCommand(cmd.GetCommand)
	rootCommand.Execute()
}
