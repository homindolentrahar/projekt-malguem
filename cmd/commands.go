package cmd

import (
	"fmt"
	"io/fs"
	stdlog "log"
	"malguem/model"
	"malguem/utils"
	"os"
	"path/filepath"
	"strings"
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

		path := "./base_template"
		malguemConfig := model.MalguemConfig{
			Name: currentDir,
			Templates: map[string]model.MalguemTemplate{
				"base_template": {
					Path: path,
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
		templateConfig := model.TemplateConfig{
			Name:     templateName,
			Language: language,
			Variables: map[string]model.TemplateConfigVariable{
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
			templateNames = append(templateNames, templateName)

			if templateInfo.Github == nil {
				continue
			}

			_, err := DownloadTemplate(templateName, templateInfo.Github.URL, templateInfo.Github.Ref, templateInfo.Github.Path)
			if err != nil {
				fmt.Printf("🌧️  Failed to download %s template from Github: %v\n", templateName, err)
				continue
			}

			fmt.Printf("🌤️  Success getting `%s` template\n", templateName)
		}

		fmt.Println("📦  Templates found: ")
		for _, name := range templateNames {
			fmt.Printf("✅  `%s`\n", name)
		}
	},
}

var MakeCommand = &cobra.Command{
	Use:   "gen",
	Short: "Generate template",
	Long:  "Command to generate template from both local and remote source",
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

		stdlog.SetFlags(0)

		// Read malguem.yaml config file
		malguemConfig, err := utils.ReadMalguemConfig()
		utils.HandleErrorExit(err)

		if _, ok := malguemConfig.Templates[templateName]; !ok {
			fmt.Printf("🌧️  %s template not found", templateName)
			return
		}

		// Read template config file on each template
		templateItem := malguemConfig.Templates[templateName]
		// Read output path of the template
		outputPath := templateItem.Output

		// Define source path of the template
		sourcePath := templateItem.Path

		if templateItem.Github != nil {
			// Download template and store it in cache
			cachePath, err := DownloadTemplate(templateName, templateItem.Github.URL, templateItem.Github.Ref, templateItem.Github.Path)
			if err != nil {
				if strings.Contains(err.Error(), "already exists") {
					err = RenderTemplate(templateName, cachePath, outputPath)
					utils.HandleErrorExit(err)

					return
				}

				fmt.Printf("🌧️  Failed to download %s template from Github: %v\n", templateName, err)
				return
			}

			sourcePath = cachePath
		}

		err = RenderTemplate(templateName, sourcePath, outputPath)
		utils.HandleErrorExit(err)
	},
}

var ListCommand = &cobra.Command{
	Use:   "list",
	Args:  cobra.NoArgs,
	Short: "List all templates",
	Long:  "List all templates available in the local machine",
	Run: func(command *cobra.Command, args []string) {
		// Get the cache dir
		cacheDir := utils.GetCacheDir()
		if cacheDir == "" {
			// Create the cache dir first
			os.MkdirAll(cacheDir, os.ModePerm)
		}

		var templates []string
		// Iterate over templates inside the cache dir
		filepath.Walk(cacheDir, func(path string, info fs.FileInfo, err error) error {
			// Check if the file is a directory
			if info.IsDir() {
				// Ensure that only list the valid template with `template.yaml` file present
				templateConfig := filepath.Join(path, "template.yaml")

				if _, err := os.Stat(templateConfig); err == nil {
					relativePath, _ := filepath.Rel(cacheDir, path)

					templates = append(templates, relativePath)

					return filepath.SkipDir
				}
			}

			return nil
		})

		fmt.Printf("🌤️  %d Templates found\n", len(templates))
		for _, template := range templates {
			fmt.Printf("   - %s\n", template)
		}
	},
}
