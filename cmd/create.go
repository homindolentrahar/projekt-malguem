package cmd

import (
	"fmt"
	"malguem/model"
	"malguem/utils"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var CreateCommad = &cobra.Command{
	Use:   "create",
	Short: "Create new template",
	Long:  "Command to create new template",
	Args:  cobra.MaximumNArgs(2),
	Run: func(command *cobra.Command, args []string) {
		var templateName string
		var language string

		// If user provide the template name by command args
		if len(args) > 0 {
			templateName = args[0]
		} else {
			// Otherwise, prompt an input
			templateName = PrompInput("Enter template name")
		}

		// If user provide the language by command args
		if len(args) > 1 {
			language = args[1]
		} else {
			// Otherwise, prompt an input
			language = PrompInput("Enter template laguage (.kt, .swift, .js)")
		}

		// Create a folder based on inputted template name
		templatePath := filepath.Join("templates", templateName)
		os.MkdirAll(templatePath, os.ModePerm)

		// Create a config file [malguem.yaml]
		configFile := "template.yaml"
		configPath := fmt.Sprintf("templates/%s/%s", templateName, configFile)
		file, err := os.Create(configPath)
		utils.HandleErrorExit(err)
		defer file.Close()

		// Write the content to the config file
		templateConfig := model.Template{
			Name:     templateName,
			Language: language,
			Output:   "./",
			Variables: map[string]model.TemplateVariable{
				"name": {
					Type:    "string",
					Default: "default",
				},
			},
		}
		templateConfigData, err := yaml.Marshal(&templateConfig)
		utils.HandleErrorExit(err)

		err = os.WriteFile(file.Name(), templateConfigData, 0644)
		utils.HandleErrorExit(err)

		// Create a package for the template
		packagePath := filepath.Join("packages", templateName+".zip")

		// Create packages folder if not exist
		os.MkdirAll("packages", os.ModePerm)

		// Create a zip file
		err = utils.CreateZipFile(templatePath, packagePath)
		utils.HandleErrorExit(err)

		fmt.Printf("🌤️  %s template created. Happy coding 🚀", templateName)
	},
}
