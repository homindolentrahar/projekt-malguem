package main

import (
	"bufio"
	"fmt"
	"malguem/model"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var rootCommand = &cobra.Command{
	Use:   "malguem",
	Short: "Boilerplate code generator",
	Long:  "A CLI tool to generate boilerplate code for multiple languages using templates.",
}
var initCommand = &cobra.Command{
	Use:   "init",
	Short: "Init malguem",
	Long:  "Init malguem in the project to use the templates",
	Run: func(cmd *cobra.Command, args []string) {
		malguemFile := "malguem.yaml"

		file, error := os.Create(malguemFile)
		if error != nil {
			panic(error)
		}
		defer file.Close()

		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		wdName := filepath.Base(wd)

		malguemConfig := model.MalguemConfig{
			Name: wdName,
			Templates: map[string]model.MalguemTemplate{
				"base_template": {
					Path: "./base_template",
				},
			},
		}

		encoder := yaml.NewEncoder(file)
		err = encoder.Encode(malguemConfig)
		if err != nil {
			panic(err)
		}

		fmt.Printf("🌤️ %s project init success", wdName)
	},
}
var createCommad = &cobra.Command{
	Use:   "create",
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
		configFile := "template.yaml"

		file, err := os.Create(fmt.Sprintf("%s/%s", configPath, configFile))
		if err != nil {
			panic(err)
		}
		defer file.Close()

		content := fmt.Sprintf(`name: %s
language: %s
outputPath: `, templateName, language)

		writer := bufio.NewWriter(file)
		_, err = writer.WriteString(content)
		if err != nil {
			panic(err)
		}

		writer.Flush()

		fmt.Printf("🌤️ %s template created\nHappy coding 🚀", templateName)
	},
}

func main() {
	rootCommand.AddCommand(initCommand)
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
