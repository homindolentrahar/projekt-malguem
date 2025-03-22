package main

import (
	"fmt"
	"malguem/cmd"
	"malguem/model"
	"malguem/utils"
	"os"
	"path/filepath"

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
		file, err := os.Create("malguem.yaml")
		utils.HandleErrorExit(err)
		defer file.Close()

		currentDir, err := utils.GetCurrentDir()
		utils.HandleErrorExit(err)

		malguemConfig := model.MalguemConfig{
			Name: currentDir,
			Templates: map[string]model.MalguemTemplate{
				"base_template": {
					Path: "./base_template",
				},
			},
		}
		data, err := yaml.Marshal(&malguemConfig)
		utils.HandleErrorExit(err)

		err = os.WriteFile(file.Name(), data, 0644)
		utils.HandleErrorExit(err)

		fmt.Printf("🌤️  %s project init success", currentDir)
	},
}

var makeCommand = &cobra.Command{
	Use:   "make",
	Short: "Make new template",
	Long:  "Command to make new template",
	Args:  cobra.MaximumNArgs(1),
	Run: func(command *cobra.Command, args []string) {
		var templateName string

		// User provide the template name by command args
		if len(args) > 0 {
			templateName = args[0]
		} else {
			// Otherwise, prompt an input
			templateName = cmd.PrompInput("Enter template name")
		}

		templateDir := filepath.Join("templates", templateName)

		// Read malguem.yaml config file
		malguemConfig, err := readMalguemConfig()
		utils.HandleErrorExit(err)

		if _, ok := malguemConfig.Templates[templateName]; !ok {
			fmt.Printf("🌧️  %s template not found", templateName)
			return
		}

		if _, err := os.Stat(templateDir); os.IsNotExist(err) {
			fmt.Printf("🌧️  %s template not found", templateName)
			return
		}
	},
}

func main() {
	rootCommand.AddCommand(initCommand)
	rootCommand.AddCommand(cmd.CreateCommad)
	rootCommand.AddCommand(makeCommand)
	rootCommand.Execute()
}

func readMalguemConfig() (*model.MalguemConfig, error) {
	file, err := os.ReadFile("malguem.yaml")
	if err != nil {
		return nil, err
	}

	var malguemConfig model.MalguemConfig
	err = yaml.Unmarshal(file, &malguemConfig)
	if err != nil {
		return nil, err
	}

	return &malguemConfig, nil
}
