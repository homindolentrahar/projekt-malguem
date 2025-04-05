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

var sectionPattern = regexp.MustCompile(`{{#([\w]+)}}([\w]+){{/[\w]+}}`)

func DownloadTemplate(template, url, ref, subDir string) (string, error) {
	homeDir, err := os.UserHomeDir()
	utils.HandleErrorReturn(err)

	cacheDir := filepath.Join(homeDir, ".malguem", "cache")
	cachePath := filepath.Join(cacheDir, template)

	// Make sure the cache directory exists
	os.MkdirAll(cacheDir, os.ModePerm)

	// Check if the template is already cached
	if _, err := os.Stat(cachePath); err == nil {
		return cachePath, fmt.Errorf("template `%s` already exists in cache", template)
	}

	// Download template and store it in cache
	tarballUrl := strings.Replace(url, "github.com", "codeload.github.com", 1) + "/tar.gz/" + ref
	err = github.CloneRepo(tarballUrl, cachePath, subDir)

	if err != nil {
		log.Error().Msg(err.Error())
		return "", err
	}

	return cachePath, nil
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

	// Make sure the output directory exists
	os.MkdirAll(output, os.ModePerm)

	// Iterate over files and folders inside path
	return filepath.Walk(path, func(pathFile string, info fs.FileInfo, err error) error {
		utils.HandleErrorReturn(err)

		// Skip the `template.yaml` file
		if filepath.Base(pathFile) == "template.yaml" {
			return nil
		}

		// Render mustache in the path
		processedPath, err := renderMustachePath(pathFile, inputData)
		utils.HandleErrorExit(err)
		// Remove the template prefix from the path
		processedPath = strings.TrimPrefix(processedPath, path)

		// Check if the `info` is a directory
		// Then make sure the directory exists and continue to the next iteration
		if info.IsDir() {
			os.MkdirAll(filepath.Join(output, processedPath), os.ModePerm)

			return nil
		}

		// Define output and render the mustache template inside content
		outputPath := filepath.Join(output, processedPath)
		err = renderMustacheContent(pathFile, outputPath, inputData)
		utils.HandleErrorExit(err)

		return nil
	})
}

func renderMustachePath(source string, data map[string]string) (string, error) {
	preprocessedTemplate := preporcessTemplate(source, data)

	result, err := mustache.Render(preprocessedTemplate, data)
	if err != nil {
		return "", err
	}

	return result, nil
}

func renderMustacheContent(source, target string, data map[string]string) error {
	file, err := os.ReadFile(source)
	utils.HandleErrorReturn(err)

	templateString := string(file)
	preprocessedTemplate := preporcessTemplate(templateString, data)

	result, err := mustache.Render(preprocessedTemplate, data)
	utils.HandleErrorReturn(err)

	return os.WriteFile(target, []byte(result), 0644)
}

func preporcessTemplate(templateString string, data map[string]string) string {
	return sectionPattern.ReplaceAllStringFunc(templateString, func(match string) string {
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
		return utils.ToTitleCase(value)
	default:
		return value
	}
}
