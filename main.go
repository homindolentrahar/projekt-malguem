package main

import (
	"bufio"
	"fmt"
	"malguem/model"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "malguem",
	Short: "Boilerplate code generator",
	Long:  "A CLI tool to generate boilerplate code for multiple languages using templates.",
}
var createCommad = &cobra.Command{
	Use:   "create [template-name]",
	Short: "Create new template",
	Long:  "Command to create new template",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var templateName string
		var language string

		// If user provide the template name by command args
		if len(args) > 0 {
			templateName = args[0]
		} else {
			// Otherwise, prompt an input
			templateName = prompInput("Enter template name")
		}

		if len(args) > 1 {
			language = args[1]
		} else {
			language = prompInput("Enter template laguage (.kt, .swift, .js)")
		}

		// Create a folder based on inputted template name
		os.MkdirAll(fmt.Sprintf("%s/%s", "templates", templateName), os.ModePerm)

		// Create a config file [malguem.yaml]
		configPath := fmt.Sprintf("templates/%s", templateName)
		configFile := "malguem.yaml"

		file, err := os.Create(fmt.Sprintf("%s/%s", configPath, configFile))
		if err != nil {
			panic(err)
		}
		defer file.Close()

		content := fmt.Sprintf(`name: %s
language: %s`, templateName, language)

		writer := bufio.NewWriter(file)
		_, err = writer.WriteString(content)
		if err != nil {
			panic(err)
		}

		writer.Flush()

		fmt.Printf("🌤️ %s template created\nHappy coding 🚀", templateName)
	},
}
var configCommand = &cobra.Command{
	Use:   "config [config-file]",
	Short: "Load config file",
	Long:  "Read and load config from .yaml file",
	Args:  cobra.MaximumNArgs(1),
	Run:   func(cmd *cobra.Command, args []string) {},
}

func main() {
	rootCommand.AddCommand(configCommand)
	rootCommand.AddCommand(createCommad)
	rootCommand.Execute()
}

func prompInput(prompt string) string {
	p := tea.NewProgram(model.PromptInputStringModel{Prompt: prompt})
	stringInput, err := p.Run()

	if err != nil {
		os.Exit(1)
	}

	value := stringInput.(model.PromptInputStringModel).Input

	if value == "" {
		return "default"
	}

	return value
}
