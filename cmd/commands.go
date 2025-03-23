package cmd

import (
	"fmt"
	stdlog "log"
	"malguem/model"
	"malguem/utils"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var InitCommand = &cobra.Command{
	Use:   "init",
	Short: "Init malguem in the project",
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

var GetCommand = &cobra.Command{
	Use:   "get",
	Short: "Get template",
	Long:  "Command to get template",
	Run: func(command *cobra.Command, args []string) {
		// Read malguem.yaml config file
		malguemConfig, err := utils.ReadMalguemConfig()
		utils.HandleErrorExit(err)

		stdlog.SetFlags(0)
		stdlog.Println("📦  Checking registered templates...")

		time.Sleep(time.Second * 1)
		var templateNames []string
		for templateName, templateInfo := range malguemConfig.Templates {
			packagePath := filepath.Join("packages", templateName+".zip")

			os.MkdirAll(filepath.Dir(packagePath), os.ModePerm)

			if _, err := os.Stat(packagePath); os.IsNotExist(err) {
				stdlog.Printf("Package is missing: %s. Getting package...", templateName)

				time.Sleep(time.Second * 1)
				// Zip or package the template from /templates folder
				err := utils.CreateZipFile(templateInfo.Path, packagePath)
				if err != nil {
					stdlog.Printf("❌ Error packaging %s: %v\n", templateName, err)
					continue
				}
			}

			templateNames = append(templateNames, templateName)
		}

		fmt.Println("🌤️ Success get all templates")
		fmt.Println("📦  Templates found: ")
		for _, name := range templateNames {
			fmt.Printf("  - %s\n", name)
		}
	},
}

var MakeCommand = &cobra.Command{
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
			templateName = PrompInput("Enter template name")
		}

		templateDir := filepath.Join("templates", templateName)

		// Read malguem.yaml config file
		malguemConfig, err := utils.ReadMalguemConfig()
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
