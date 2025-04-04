package cmd

import (
	"fmt"
	"io/fs"
	"malguem/utils"
	"malguem/utils/git/github"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cbroglie/mustache"
	"github.com/rs/zerolog/log"
)

var homeDir, err = os.UserHomeDir()

var sectionPattern = regexp.MustCompile(`{{#([\w]+)}}([\w]+){{/[\w]+}}`)

func DownloadTemplate(template, url, ref, subDir string) (error, string) {
	homeDir, err := os.UserHomeDir()
	utils.HandleErrorReturn(err)

	cacheDir := filepath.Join(homeDir, ".malguem", "cache")
	cachePath := filepath.Join(cacheDir, template)

	// Make sure the cache directory exists
	os.MkdirAll(cacheDir, os.ModePerm)

	// Check if the template is already cached
	if _, err := os.Stat(cachePath); err == nil {
		return fmt.Errorf("Template `%s` already exists in cache\n", template), cachePath
	}

	// Download template and store it in cache
	tarballUrl := strings.Replace(url, "github.com", "codeload.github.com", 1) + "/tar.gz/" + ref
	err = github.CloneRepo(tarballUrl, cachePath, subDir)

	if err != nil {
		log.Error().Msg(err.Error())
		return err, ""
	}

	return nil, cachePath
}

func RenderTemplate(template, path, output string) error {
	// Read template config `template.yaml`
	templateConfigPath := filepath.Join(path, "template.yaml")
	templateConfig, err := utils.ReadTemplate(templateConfigPath)
	utils.HandleErrorReturn(err)

	// Read template variables
	var inputData = make(map[string]string)
	for key := range templateConfig.Variables {
		variable := templateConfig.Variables[key]
		inputPrompt := PrompInput(fmt.Sprintf("%s (%v)", variable.Prompt, variable.Default))

		// If user input is empty, use default value
		if inputPrompt == "" {
			inputData[key] = variable.Default.(string)
		} else {
			inputData[key] = inputPrompt
		}
	}

	return filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		utils.HandleErrorReturn(err)

		// Ignore directories and `template.yaml` file
		if info.IsDir() || filepath.Base(path) == "template.yaml" {
			return nil
		}

		outputPath := filepath.Join(output, filepath.Base(path))
		err = renderMustache(path, outputPath, inputData)
		utils.HandleErrorExit(err)

		return nil
	})
}

func renderMustache(source, target string, data map[string]string) error {
	file, err := os.ReadFile(source)
	utils.HandleErrorReturn(err)

	templateString := string(file)
	preprocessedTemplate := sectionPattern.ReplaceAllStringFunc(templateString, func(match string) string {
		// Extract format type and variable name
		// {{format}} varName {{/format}}
		//     1         2         3
		matches := sectionPattern.FindStringSubmatch(match)
		// Invalid format
		if len(matches) < 3 {
			return match
		}

		format, variable := matches[1], matches[2]

		// Check if variable exists inside userData
		value, exists := data[variable]
		if !exists {
			return match
		}

		return applyFormat(format, value)
	})

	result, err := mustache.Render(preprocessedTemplate, data)
	utils.HandleErrorReturn(err)

	return os.WriteFile(target, []byte(result), 0644)
}

func applyFormat(format, value string) string {
	switch format {
	case "pascal_case":
		return utils.ToPascalCase(value)
	case "camel_case":
		return utils.ToCamelCase(value)
	case "snake_case":
		return utils.ToSnakeCase(value)
	case "kebab_case":
		return utils.ToKebabCase(value)
	case "uppercase":
		return strings.ToUpper(value)
	case "lowercase":
		return strings.ToLower(value)
	case "titlecase":
		return strings.Title(value)
	default:
		return value
	}
}
